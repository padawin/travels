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
			.when("/picture/:travelId/:place/:picture", {
				templateUrl: "partials/picture.html",
				controller: "PictureCtrl"
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
				$location.path('/pictures/' + $scope.travelId).replace();
				return;
			}

			pics.sort();
			places.sort();

			for (var p in places) {
				$scope.places[p] = {
					'name': places[p],
					'preview': pics[p][0|Math.random()*pics[p].length]
				};
			}

			$scope.title =  $scope.travels[$scope.travelId].title;
		});
	});

	travelsApp.controller('PicturesListCtrl', function($rootScope, $scope, $http, $routeParams, $location){
		var travelId = $routeParams.travelId,
			place = $routeParams.place,
			picture, file;

		$rootScope.$emit('display-places-list', 1);

		getTravels($scope, $http, function($scope) {
			if (!$scope.travels[travelId]) {
				$location.url('/');
				return;
			}

			var pictures, placeIndex, subtitle,
				backLink = '#/places/' + travelId;
			if (!place) {
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
				subtitle = place;
				placeIndex = $scope.travels[travelId].places.indexOf(place);
				if (!~placeIndex) {
					$location.url('/');
					return;
				}

				pictures = $scope.travels[travelId].pics[placeIndex];
			}

			pictures.sort();
			$scope.backLink = backLink;
			$scope.title =  $scope.travels[travelId].title;
			$scope.subtitle = subtitle;
			$scope.travelId = travelId;
			$scope.place = place;
			$scope.pictures = pictures;
		});
	});

	travelsApp.controller('PictureCtrl', function($rootScope, $scope, $http, $routeParams, $location){
		var travelId = $routeParams.travelId,
			place = $routeParams.place,
			picture = parseInt($routeParams.picture);

		$rootScope.$emit('display-places-list', 1);

		getTravels($scope, $http, function($scope) {
			if (!$scope.travels[travelId]) {
				$location.url('/');
				return;
			}

			var placeIndex,
				next = null,
				prev = null,
				backLinkTravel = '#/places/' + travelId,
				backLinkPlace = '',
				pictures;

			$scope.title = $scope.travels[travelId].title;
			if (!place) {
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

				$scope.subtitle = place;
				pictures = $scope.travels[travelId].pics[placeIndex];
				backLinkPlace += '#/pictures/' + travelId + '/' + place;
			}

			if (picture > 0) {
				prev = picture - 1;
			}
			if (picture < pictures.length - 1) {
				next = picture + 1;
			}
			picture = pictures[picture];
			$scope.backLinkTravel = backLinkTravel;
			$scope.backLinkPlace = backLinkPlace;
			$scope.travelId = travelId;
			$scope.place = place;
			$scope.picture = picture;
			$scope.next = next;
			$scope.prev = prev;
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

				pics.sort();
				places.sort();

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
