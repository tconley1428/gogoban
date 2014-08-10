var BOARD_SIZE = 19

var GoGoBan = angular.module('GoGoBan', []);

GoGoBan.controller('GoGoBanCtrl', function ($scope) {
  $scope.board = [];
  for (var i = 0; i < BOARD_SIZE; i++) {
	  var row = [];
	  for (var j = 0; j < BOARD_SIZE; j++) {
	  	row.push("+");
	  };
	  $scope.board.push(row)
  }
  $scope.print = function(x,y){
  	alert("("+x+","+y+")")
  }
});