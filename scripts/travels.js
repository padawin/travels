(function(){
var travelsApp = angular.module('travelsApp', []);

	travelsApp.controller('TravelsListCtrl', function($scope, $http){
		$http.get('data/travels.json').success(function(data) {
			var t, pics;
			for (t in data) {
				data[t].previewPics = data[t].pics
					// duplicate initial pics list
					.slice(0)
					// randomize
					.sort(function() { return 0.5 - Math.random() })
					// takes the 4 first
					.splice(0, 4);
			}

			$scope.travels = data;
			$scope.orderProp = 'title';
		});
	});

	window.travelsApp = travelsApp;
})();
