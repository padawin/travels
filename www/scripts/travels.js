(function(){
	var travelsApp,
		travelsAppMenu,
		travels,
		PlacesListCtrl,
		getTravels;

	travelsApp = angular.module("travelsApp", ['ngRoute'], function($routeProvider, $locationProvider) {
		$locationProvider.hashPrefix('');
		$routeProvider
			.when("/", {
				templateUrl: "partials/travels-list.html",
				controller: "TravelsListCtrl"
			})
			.when("/places/:travelId", {
				templateUrl: "partials/places-list.html",
				controller: "PlacesListCtrl"
			})
			.when("/pictures/:travelId", {
				templateUrl: "partials/pictures-list.html",
				controller: "PicturesListCtrl"
			})
			.when("/pictures/:travelId/:place", {
				templateUrl: "partials/pictures-list.html",
				controller: "PicturesListCtrl"
			})
			.otherwise({redirectTo: "" });
	});

	getTravels = function($scope, $http, doneCallback)
	{
		if (travels != null) {
			$scope.travels = travels;
			doneCallback && doneCallback($scope);
		}
		else {
			$http.get('data/travels.json').success(function(data) {
				travels = data;
				$scope.travels = data;
				doneCallback && doneCallback($scope);
			});
		}
	};

	travelsApp.controller('TravelsListCtrl', function($rootScope, $scope, $http){
		$rootScope.$emit('display-places-list', 0);
		$scope.orderProp = 'title';
		getTravels($scope, $http, function($scope){
			var t, pics;
			for (t in $scope.travels) {
				$scope.travels[t].id = t;
				$scope.travels[t].previewPics = [].concat.apply([], $scope.travels[t].pics)
					// duplicate initial pics list
					.slice(0)
					// randomize
					.sort(function() { return 0.5 - Math.random() })
					// takes the 4 first
					.splice(0, 4);
			}
		});
	});

	travelsApp.controller('PlacesListCtrl', function($rootScope, $scope, $http, $routeParams, $location){
		$rootScope.$emit('display-places-list', 0);
		$scope.travelId = $routeParams.travelId;
		getTravels($scope, $http, function($scope){
			$scope.places = [];
			var p, places = $scope.travels[$scope.travelId].places,
				pics = $scope.travels[$scope.travelId].pics;

			if (places.length == 0) {
				$location.url('/pictures/' + $scope.travelId);
				return;
			}

			for (var p in places) {
				$scope.places[p] = {
					'name': places[p],
					'preview': pics[p][0|Math.random()*pics[p].length]
				};
			}
		});
	});

	travelsApp.controller('PicturesListCtrl', function($rootScope, $scope, $http, $routeParams, $location){
		var travelId = $routeParams.travelId,
			place = $routeParams.place;

		$rootScope.$emit('display-places-list', 1);

		getTravels($scope, $http, function($scope) {
			if (!$scope.travels[travelId]) {
				$location.url('/');
				return;
			}

			var pictures, placeIndex;
			if (!place) {
				$scope.title = $scope.travels[travelId].title;
				// 0 place for this travel
				if ($scope.travels[travelId].places.length == 0) {
					pictures = $scope.travels[travelId].pics;
				}
				// else get all the pictures of the travel
				else {
					pictures = [].concat.apply([], $scope.travels[travelId].pics);
				}
			}
			else {
				placeIndex = $scope.travels[travelId].places.indexOf(place);
				if (!~placeIndex) {
					$location.url('/');
					return;
				}

				$scope.title = place;
				pictures = $scope.travels[travelId].pics[placeIndex];
			}

			$scope.pictures = pictures;
		});
	});

	PlacesListMenuCtrl = function($rootScope, $scope, $http, $routeParams){
		$rootScope.$on('display-places-list', function(e, display) {
			if (!display) {
				$scope.menuPlaces = [];
				return;
			}

			$scope.travelId = $routeParams.travelId;

			getTravels($scope, $http, function($scope){
				$scope.menuPlaces = [];
				var p, places = $scope.travels[$scope.travelId].places,
					pics = $scope.travels[$scope.travelId].pics;

				for (var p in places) {
					$scope.menuPlaces[p] = {
						'name': places[p],
						'preview': pics[p][0|Math.random()*pics[p].length]
					};
				}
			});
		});
	};

	window.travelsApp = travelsApp;
	window.PlacesListMenuCtrl = PlacesListMenuCtrl;
})();