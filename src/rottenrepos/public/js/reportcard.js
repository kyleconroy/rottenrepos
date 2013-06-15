$(document).ready(function() {
  var header = $("h2");
  var user = header.data('user');
  var repo = header.data('repo');

  var url = "/github/" + user + "/" + repo;

  var jqxhr = $.ajax({
    url: url,
    headers: {"X-PJAX": "True"}
  }).done(function(e) { 
    console.log(e);
    $("#card").html(e);
  });
});
