(function() {
  var inputRemote, paperToast, lyricsElement;

  HTMLImports.whenReady(function() {
    inputRemote = document.querySelector('#input-remote');
    paperToast = document.querySelector('paper-toast');
    lyricsElement = document.querySelector('#lyrics');
    document.addEventListener('autocomplete-selected', onSelect);
    document.addEventListener('autocomplete-change', onChange)
    inputRemote.querySelector("input").addEventListener('change', function(e) {
      if (document.activeElement == this) {
        if (inputRemote.querySelector("paper-item[aria-selected=true]") != null) {
          return false;
        }
        inputRemote._fireEvent({
          text: inputRemote.text,
          value: inputRemote.text
        }, 'selected');
        inputRemote.reset();
        return true;
      }
    });
    inputRemote.$.autocompleteInput.focus();
    lyrics.style.setProperty('display', 'none');
  });

  function onSelect(event) {
    var selected = event.detail.text;
    paperToast.text = 'Search lyrics for: ' + selected;
    paperToast.show();
    searchLyrics(selected);
  }

  function searchLyrics(lyrics) {
    var url = '/lyrics?lyrics=' + lyrics;
    var req = new XMLHttpRequest();
    req.open('GET', encodeURI(url));
    req.onload = function() {
      paperToast.hide();
      if (req.status === 200) {
        var data = JSON.parse(req.response);
        lyricsElement.style.removeProperty('display');
        lyricsElement.querySelector("#lyrics-title").innerHTML = data.lyrics.Title;
        lyricsElement.querySelector("#lyrics-lyrics").innerHTML = data.lyrics.Lyrics;
      }
    };
    req.send();
  }

  function onChange(event) {
    var search = event.detail.option.text;
    var url = '/suggest?lyrics=' + search;
    var req = new XMLHttpRequest();
    req.open('GET', encodeURI(url));
    req.onload = function() {
      if (req.status === 200) {
        var data = JSON.parse(req.response);
        var arr = data.suggestions.map(function(obj) {
          return {
            text: obj,
            value: obj
          };
        });
        inputRemote.suggestions(arr);
      }
    };
    req.send();
  }
})();