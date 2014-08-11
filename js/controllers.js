var BOARD_SIZE = 19

var GoGoBan = angular.module('GoGoBan', []);

GoGoBan.controller('GoGoBanCtrl', function ($scope) {

  var connection = new WebSocket('ws://'+window.location.host+'/ws/game', []);

  connection.onmessage = function(e){
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
  	if($scope.board[x][y].Player == "empty"){
  		connection.send(JSON.stringify({"X":x,"Y":y,"Player":"empty"}))
  	}
  }
});

GoGoBan.controller('GoGoBanLobbyCtrl', function ($scope) {
  $scope.players = []
  $scope.name = ''
  $scope.opponent = ''
  var connection = new WebSocket('ws://'+window.location.host+'/ws/lobby', []);
  connection.onmessage = function(e){
  	obj = JSON.parse(e.data)
  	if(obj.Source != undefined){
  		$scope.opponent = obj.Source
  		$('#modal').modal({})

  	}else{
	  	$scope.players = obj
  	}
	$scope.$apply()
  }

  $scope.sendName = function(name){
  	$scope.name = name
  	connection.send(name)
  }
  $scope.request = function(name){
	$scope.opponent = name
	$scope.$apply()
	$('#request').modal({})
  	connection.send(JSON.stringify({"Target":name,"Source":$scope.name}))
  }
});