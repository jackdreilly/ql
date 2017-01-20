function selectAlternative(e) {
  console.log(e);
}
(function() {
  var inputRemote, paperToast, spinner, lyricsCard, additionalResults;

  HTMLImports.whenReady(function() {
    inputRemote = document.querySelector('#input-remote');
    paperToast = document.querySelector('paper-toast');
    spinner = document.querySelector('paper-spinner');
    lyricsCard = document.querySelector('#lyrics');
    lyricsItem = lyricsCard.querySelector('#lyrics');
    additionalResults = document.querySelector('iron-list');
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
    spinner.active = true;
    req.onload = function() {
      paperToast.hide();
      spinner.active = false;
      if (req.status === 200) {
        var data = JSON.parse(req.response);
        lyricsCard.heading = data.Lyrics.Title;
        lyricsItem.querySelector('.card-content').innerHTML = data.Lyrics.Lyrics;
        additionalResults.items = data.Alternatives;
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