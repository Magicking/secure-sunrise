var endpoint = '/api';

setInterval(function(){
$.ajax({url: endpoint + "/feeds?name=TODO",})
  .done(function( data ) {
    var url = data[Math.floor(Math.random() * data.length)];
    $("#image0").attr("src", url);
    console.log(url);
   timeout: 1000 //in milliseconds
  });
}, 10000);
