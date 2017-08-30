var endpoint = '/api';

setInterval(function(){
$.ajax({url: endpoint + "/feeds?name=TODO",})
  .done(function( data ) {
    $("#image0").attr("src", data[0]);
   timeout: 1000 //in milliseconds
});
  });
}, 10000);
