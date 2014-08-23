'use strict';

var getCookie = function (cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1);
        if (c.indexOf(name) != -1) return c.substring(name.length,c.length);
    }
    return "";
}

var GoGoBan = angular.module('GoGoBan', [
	'ngRoute',
	'GoGoBanControllers'
])



GoGoBan.service('login', function(){
    this.isUserLoggedIn = function (){
      return getCookie("username")!=""
    };

    this.User = function (){
      return getCookie("username")
    }
});

GoGoBan.run(function($rootScope, $location,login) {
  $rootScope.$on('$routeChangeStart', function () {
    if (!login.isUserLoggedIn()) {
      $location.path('/login');
    }
  })
});

GoGoBan.
  config(['$routeProvider', function($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: '/www/partials/login.html',
        controller: 'LoginCtrl'
      })
      .when('/expired', {
        templateUrl: '/www/partials/expired.html',
        controller: 'ExpiredCtrl'
      })
      .when('/login', {
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