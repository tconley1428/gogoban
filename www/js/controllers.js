var BOARD_SIZE = 19

var GoGoBanControllers = angular.module('GoGoBanControllers', []);

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1);
        if (c.indexOf(name) != -1) return c.substring(name.length,c.length);
    }
    return "";
}

GoGoBanControllers.controller('BoardCtrl', function ($scope) {
  $scope.room = window.location.hash.substring(window.location.hash.lastIndexOf('/')+1)
  $scope.IsTurn = false
  var connection = new WebSocket('wss://'+window.location.host+'/wss/game/'+$scope.room, []);

  connection.onmessage = function(e){
    if(e.data=="Your Turn"){
      $scope.IsTurn = true
      return
    }
    obj = JSON.parse(e.data)

    $scope.board[obj.X][obj.Y] = {"Player":obj.Player}
    $scope.$apply()
  }

  $scope.board = [];
  for (var i = 0; i < BOARD_SIZE; i++) {
    var row = [];
    for (var j = 0; j < BOARD_SIZE; j++) {
      row.push({"Player":"empty"});
    };
    $scope.board.push(row)
  }
  $scope.click = function(x,y){
    if($scope.board[x][y].Player == "empty" && $scope.IsTurn){
      connection.send(JSON.stringify({"X":x,"Y":y}))
      $scope.IsTurn=false
    }
  }
});

GoGoBanControllers.controller('LoginCtrl', function ($scope) {

})

GoGoBanControllers.controller('MenuCtrl', function ($scope) {


  $scope.username = getCookie("username")
})

GoGoBanControllers.controller('LobbyCtrl', function ($scope,$location) {
   var connection = new WebSocket('wss://'+window.location.host+'/wss/lobby', []);

  $scope.players = []

  $scope.sendName = function(name){
    $scope.name = name
    connection.send(name)
  }

  connection.onopen= function(e){
      $scope.sendName(getCookie("username"))
  }

  function isArray(what) {
    return Object.prototype.toString.call(what) === '[object Array]';
  }

  $scope.opponent = ''
  connection.onmessage = function(e){
    obj = JSON.parse(e.data)
    if(obj.Status == "request"){
      $scope.opponent = obj.Source
      $('#challenge').modal({})
    }else if(obj.SessionID !=undefined){
      $('.modal').modal('hide')
      $location.path('/board/'+obj.SessionID)
    }else if(obj.Status===undefined && isArray(obj)){
      $scope.players = obj
    }
  $scope.$apply()
  }


  $scope.request = function(name){
  $scope.opponent = name
  $('#request').modal({})
    connection.send(JSON.stringify({"Status":"request","Target":name,"Source":$scope.name}))
  }

  $scope.acceptRequest = function(){
    connection.send(JSON.stringify({"Status":"accept","Target":$scope.opponent,"Source":$scope.name}))
  }
  $scope.declineRequest = function(){
    connection.send(JSON.stringify({"Accept":false}))
  }
  $scope.cancelRequest = function(){
    connection.send(JSON.stringify({"Cancel":true}))
  }
});