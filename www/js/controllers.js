var BOARD_SIZE = 19

var GoGoBanControllers = angular.module('GoGoBanControllers', []);

function getQueryParams(qs) {
    qs = qs.split("+").join(" ");

    var params = {}, tokens,
        re = /[?&]?([^=]+)=([^&]*)/g;

    while (tokens = re.exec(qs)) {
        params[decodeURIComponent(tokens[1])]
            = decodeURIComponent(tokens[2]);
    }

    return params;
}
GoGoBanControllers.controller('ExpiredCtrl', function ($scope, $location) {
  function delete_cookie( name ) {
    document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
  }
  delete_cookie("session")
  delete_cookie("username")
  $location.path("/login")
})

GoGoBanControllers.controller('BoardCtrl', function ($scope, $location) {
  $scope.room = window.location.hash.substring(window.location.hash.lastIndexOf('/')+1)
  $scope.IsTurn = false
  var connection = new WebSocket('wss://'+window.location.host+'/wss/game/'+$scope.room, []);

  connection.onmessage = function(e){
    if(e.data=="Your Turn"){
      $scope.IsTurn = true
      return
    } else if (e.data=="InvalidMove"){
      return
    } else if (e.data=="InvalidLogin"){
      $location.path("/expired")
      $scope.$apply()
    }
    obj = JSON.parse(e.data)

    $scope.board[obj.Loc.X][obj.Loc.Y] = {"Player":obj.Player}
    $scope.$apply()
  }
  connection.onclose=function(e){
    $location.path("/expired")
    $scope.$apply()
  }
  connection.onerror=function(e){
    $location.path("/expired")
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
      connection.send(JSON.stringify({"Loc":{"X":x,"Y":y}}))
      $scope.IsTurn=false
    }
  }
});

GoGoBanControllers.controller('LoginCtrl', function ($scope,$routeParams,$location) {
  if (getCookie("username")!="") {
    $location.path('/lobby');
  }
  $scope.failed = $routeParams["failed"]!=undefined
  $scope.email = $routeParams["failed"]
  $scope.register = function(){
    $("#register").modal({})
  }
})

GoGoBanControllers.controller('MenuCtrl', function ($scope,login) {
  $scope.username = login.User()
})

GoGoBanControllers.controller('LobbyCtrl', function ($scope,$location,login) {
  var connection = new WebSocket('wss://'+window.location.host+'/wss/lobby', []);

  $scope.players = []

  connection.onclose=function(e){
    $location.path("/expired")
    $scope.$apply()
  }
  connection.onerror=function(e){
    $location.path("/expired")
    $scope.$apply()
  }



  function isArray(what) {
    return Object.prototype.toString.call(what) === '[object Array]';
  }

  $scope.opponent = ''
  connection.onmessage = function(e){
    if(e.data==="InvalidLogin"){
      $location.path("/expired")
      $scope.$apply()
      return
    }
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
    connection.send(JSON.stringify({"Status":"request","Target":name}))
  }

  $scope.acceptRequest = function(){
    connection.send(JSON.stringify({"Status":"accept","Target":$scope.opponent}))
  }
  $scope.declineRequest = function(){
    connection.send(JSON.stringify({"Accept":false}))
  }
  $scope.cancelRequest = function(){
    connection.send(JSON.stringify({"Cancel":true}))
  }
});