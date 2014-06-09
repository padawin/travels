(function(){
	var travelsApp = angular.module('travelsApp', []),
		travels,
		getTravels;

	getTravels = function($scope, $http)
	{
		if (travels != null) {
			$scope.travels = travels;
		}
		else {
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
		}
	};


	travelsApp.controller('TravelsListCtrl', function($scope, $http){
		getTravels($scope, $http);
	});

	window.travelsApp = travelsApp;
})();
