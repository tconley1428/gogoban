'use strict';

var GoGoBan = angular.module('GoGoBan', [
	'ngRoute',
	'GoGoBanControllers'
]);

GoGoBan.
  config(['$routeProvider', function($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: '/www/partials/login.html',
        controller: 'LoginCtrl'
      })
      .when('/lobby', {
        templateUrl: '/www/partials/lobby.html',
        controller: 'LobbyCtrl'
      })
      .when('/board/:sessionid', {
        templateUrl: '/www/partials/board.html',
        controller: 'BoardCtrl'
      });
  }]);