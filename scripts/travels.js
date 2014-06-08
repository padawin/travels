(function(){
var travelsApp = angular.module('travelsApp', []);

	travelsApp.controller('TravelsListCtrl', function($scope){
		$scope.travels = [
			{
				'id': '1',
				'title': 'Portugal 2011',
				'pics': ['myPic1.jpg', 'myPic2.jpg', 'myPic3.jpg', 'myPic4.jpg']
			}
		];
	});

	window.travelsApp = travelsApp;
})();
