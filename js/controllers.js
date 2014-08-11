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

GoGoBan.controller('GoGoBanLobbyCtrl', function ($scope,$window) {

  $scope.players = []
  $scope.name = ''
  $scope.opponent = ''
  var connection = new WebSocket('ws://'+window.location.host+'/ws/lobby', []);
  connection.onmessage = function(e){
  	obj = JSON.parse(e.data)
  	if(obj.Status == "request"){
  		$scope.opponent = obj.Source
  		$('#modal').modal({})
	}else if(obj.Status == "accept"){
		$window.location.href = $window.location.host +'/board'
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