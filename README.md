# [Deplecated] lack of base concept

# goplate
# Thanks for libraries
https://github.com/PuerkitoBio/goquery<br>
https://github.com/yosssi/gohtml<br>
https://github.com/ditashi/jsbeautifier-go/jsbeautifier<br>
https://github.com/gorilla/css/scanner<br>
https://github.com/revel/revel<br>
<br>
# / : template style Front-end preprocessor for angularjs
- Concept
1) embeded css<br>
css block only affect within current plate block(if inherit=true, subplate blocks are also affeted by current css block)<br>
<br>
2) structured css<br>
css block support structured css<br>
<br>
3) embeded javascript(intergrated into angularjs)<br>
script block will be change to angularjs event, and parents of script block will be controller<br>
attribute inject of script will inject module into current Controller<br>
<br>
4) each plate turn to html element<br>
if you have plate has name subject, when you use subject element, it will be replaced with plate<br>
all other inline html, script plate element have will override original plate<br>
<br>
- Current<br>
embeded css<br>
structured css<br>
embeded javascript(intergrated into angularjs)<br>
plate trun to html element<br>
<br>
- Future<br>
less or scss support(or both with type attribute)<br>
dependency management<br>
holder & place in plate<br>
<br>
- Bug<br>
DOM that is inserted to plate's child does not secure order<br>
Recursive will be shutdown program<br>
<br>
## Before Compile
item.html
```html
<plate name="item">
	<li>
		Product Name : {{item.name}}
		Count : {{item.count}}
	</li>
</plate>
```
itemlist.html
```html
<plate name="itemList">
	<css>
	</css>
	<ul>
		- Product List
		<item ng-repeat="item in itemList">
			<script event="click(item)">
			$scope.$emit("item.click", item);
			</script>
		</item>
		<script>
		$scope.itemList = [];
		</script>
	</ul>
</plate>
```
shop.html
```html
<plate name="shop">
	<div>
		<itemList>
				<script>
				$scope.itemList = $scope.$parent.itemList;
				$scope.$on("item.click", function(e, item) {
					$scope.$emit("buy", item);
				});
				</script>
		</itemList>
		Cart : {{ money }}
		<itemList>
				<script>
				$scope.itemList = $scope.$parent.cartList;
				$scope.$on("item.click", function(e, item) {
					$scope.$emit("sell", item);
				});
				</script>
		</itemList>
		<script>
			var itemHash = {
				1 : {name:"A", cost:500},
				2 : {name:"B", cost:900}
			};
			$scope.getItemInfo = function(id) {
				return itemHash[id];
			}
			$scope.money = 5000;
		</script>
		<script>
			$scope.itemInstList = [];
			$scope.cartInstList = [];
			$scope.itemInstList.push({id:1, count:5});
			$scope.itemInstList.push({id:2, count:5});
		</script>
		<script>
			function updateList(itemInstList, itemList) {
				itemList.splice(0, itemList.length);
				for(var i=0; i<itemInstList.length; ++i) {
					var itemInst = itemInstList[i];
					itemList.push({
						id : itemInst.id,
						name : $scope.getItemInfo(itemInst.id).name,
						count : itemInst.count,
						_inst : itemInst
					});
				}
			}
			$scope.itemList = [];
			$scope.cartList = [];
			$scope.updateItemList = function() {
				updateList($scope.itemInstList, $scope.itemList);
			};
			$scope.updateCartList = function() {
				updateList($scope.cartInstList, $scope.cartList);
			};
			$scope.updateItemList();
			$scope.updateCartList();
		</script>
		<script>
			function transaction(args) {
				var fromList = args.from;
				var toList = args.to;
				var item = args.item;
				var selectFunc = args.selectFunc;
				var validFunc = args.validFunc;
				var isValid = true;
				if(validFunc != null) {
					isValid = validFunc(item);
				}
				if(isValid) {
					item.count--;
					if(item.count <= 0) {
						for(var i=0; i<fromList.length; ++i) {
							if(fromList[i] == item) {
								fromList.splice(i, 1);
								break;
							}
						}
					}
					var cartInst = null;
					for(var i=0; i<toList.length; ++i) {
						if(toList[i].id == item.id) {
							var isSelect = true;
							if(selectFunc != null) {
								isSelect = selectFunc(toList[i]);
							}
							if(isSelect) {
								cartInst = toList[i];
								break;
							}
						}
					}
					if(cartInst == null) {
						cartInst = {id:item.id, count:0};
						toList.push(cartInst);
					}
					cartInst.count++;
				}
			}
			$scope.$on("buy", function(e, item) {
				transaction({
					from : $scope.itemInstList,
					to : $scope.cartInstList,
					item : item._inst,
					selectFunc : function(it) {
						return it.count < 3;
					},
					validFunc : function(it) {
						var itemInfo = $scope.getItemInfo(it.id);
						if(itemInfo.cost <= $scope.money) {
							$scope.money -= itemInfo.cost;
							return true;
						} else {
							return false;
						}
					}
				});
				$scope.updateItemList();
				$scope.updateCartList();
			});
			$scope.$on("sell", function(e, item) {
				transaction({
					from : $scope.cartInstList,
					to : $scope.itemInstList,
					item : item._inst,
					selectFunc : function(it) {
						return it.count < 2;
					},
					validFunc : function(it) {
						var itemInfo = $scope.getItemInfo(it.id);
						$scope.money += Math.round(itemInfo.cost * 0.5);
						return true;
					}
				});
				$scope.updateItemList();
				$scope.updateCartList();
			});
		</script>
	</div>
</plate>
```
index.html
```html
<html>
<head>
	<meta charset="utf8" />
</head>
<body>
<shop>
</shop>
</body>
<script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.3.14/angular.min.js"></script>
</html>
```

## After Compile
```HTML
<html>
	<head>
		<meta charset="utf8"/>
		<style>
			  .ABC.genclass_1_1 {
			
			background-color: black;
			}
			
			  .ABC {
			
			background-color: red;
			
			
			
			}
			  .ABC .BCD.genclass_1_2 {
			
			background-color: orange;
			}
			  .ABC div.genclass_1_3 {
			
			background-color: green;
			}
			  .ABC:hover {
			
			background-color: blue;
			}
		</style>
	</head>
	<body ng-app="myApp">
		<div ng-controller="Ctrl_4">
			<ul ng-controller="Ctrl_2">
				<li class="ABC  genclass_1_1" ng-repeat="subject in subjectList" ng-if="subjectList.length &gt; 0" ng-controller="Ctrl_1" ng-click="EventHandler_1_1(subject)" ng-mousemove="EventHandler_1_2(subject)">
					<span class="BCD  genclass_1_2" sid="{{subject.Sid}}">
						{{subject.Name}}
						<div ef="123" class="genclass_1_3">
							EF
						</div>
					</span>
				</li>
			</ul>
			<div ng-controller="Ctrl_3">
				<span>
					{{Sid}}
				</span>
			</div>
		</div>
		<div ng-controller="Ctrl_9">
			<ul ng-controller="Ctrl_6">
						- Product List
						
				<li ng-repeat="item in itemList" ng-controller="Ctrl_5" ng-click="EventHandler_5_1(item)">
							Product Name : {{item.name}}
							Count : {{item.count}}
						
				</li>
			</ul>
					Cart : {{ money }}
					
			<ul ng-controller="Ctrl_8">
						- Product List
						
				<li ng-repeat="item in itemList" ng-controller="Ctrl_7" ng-click="EventHandler_7_1(item)">
							Product Name : {{item.name}}
							Count : {{item.count}}
						
				</li>
			</ul>
		</div>
		<script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.3.14/angular.min.js">
		</script>
		<script>
			var myApp = angular.module('myApp',[]);myApp.controller('Ctrl_2', ['$scope', '$rootScope', function($scope, $rootScope) {
			    //Load Event
			    /*
					$scope.arg[0];
					$scope.arg[1];
					$scope.arg[2];
					*/
			    $scope.subjectList = [{
			        Name: "TestName",
			        Sid: "123"
			    }, {
			        Name: "TestName2",
			        Sid: "1232"
			    }, {
			        Name: "TestName3",
			        Sid: "1233"
			    }];
			}]);
			myApp.controller('Ctrl_4', ['$scope', '$rootScope', function($scope, $rootScope) {
			    $scope.$on('subject.select', function(e, Sid) {
			        if ($scope == e.targetScope) {
			            return;
			        }
			        console.log($scope, e);
			        $scope.$broadcast('subject.select', Sid);
			    });
			}]);
			myApp.controller('Ctrl_5', ['$scope', '$rootScope', function($scope, $rootScope) {
			    $scope.EventHandler_5_1 = function(item) {
			        $scope.$emit("item.click", item);
			    }
			}]);
			myApp.controller('Ctrl_7', ['$scope', '$rootScope', function($scope, $rootScope) {
			    $scope.EventHandler_7_1 = function(item) {
			        $scope.$emit("item.click", item);
			    }
			}]);
			myApp.controller('Ctrl_8', ['$scope', '$rootScope', function($scope, $rootScope) {
			    $scope.itemList = [];
			
			    $scope.itemList = $scope.$parent.cartList;
			    $scope.$on("item.click", function(e, item) {
			        $scope.$emit("sell", item);
			    });
			}]);
			myApp.controller('Ctrl_1', ['$scope', '$rootScope', '$http', '$q', function($scope, $rootScope, $http, $q) {
			    $scope.EventHandler_1_1 = function(subject) {
			        $scope.$emit('subject.select', subject.Sid);
			        //alert(subject.Sid);
			    }
			
			    $scope.EventHandler_1_2 = function(subject) {
			        console.log("B");
			    }
			}]);
			myApp.controller('Ctrl_3', ['$scope', '$rootScope', function($scope, $rootScope) {
			    $scope.$on('subject.select', function(e, Sid) {
			        $scope.Sid = Sid;
			    });
			}]);
			myApp.controller('Ctrl_6', ['$scope', '$rootScope', function($scope, $rootScope) {
			    $scope.itemList = [];
			
			    $scope.itemList = $scope.$parent.itemList;
			    $scope.$on("item.click", function(e, item) {
			        $scope.$emit("buy", item);
			    });
			}]);
			myApp.controller('Ctrl_9', ['$scope', '$rootScope', function($scope, $rootScope) {
			    var itemHash = {
			        1: {
			            name: "A",
			            cost: 500
			        },
			        2: {
			            name: "B",
			            cost: 900
			        }
			    };
			    $scope.getItemInfo = function(id) {
			        return itemHash[id];
			    }
			    $scope.money = 5000;
			
			    $scope.itemInstList = [];
			    $scope.cartInstList = [];
			    $scope.itemInstList.push({
			        id: 1,
			        count: 5
			    });
			    $scope.itemInstList.push({
			        id: 2,
			        count: 5
			    });
			
			    function updateList(itemInstList, itemList) {
			        itemList.splice(0, itemList.length);
			        for (var i = 0; i < itemInstList.length; ++i) {
			            var itemInst = itemInstList[i];
			            itemList.push({
			                id: itemInst.id,
			                name: $scope.getItemInfo(itemInst.id).name,
			                count: itemInst.count,
			                _inst: itemInst
			            });
			        }
			    }
			    $scope.itemList = [];
			    $scope.cartList = [];
			    $scope.updateItemList = function() {
			        updateList($scope.itemInstList, $scope.itemList);
			    };
			    $scope.updateCartList = function() {
			        updateList($scope.cartInstList, $scope.cartList);
			    };
			    $scope.updateItemList();
			    $scope.updateCartList();
			
			    function transaction(args) {
			        var fromList = args.from;
			        var toList = args.to;
			        var item = args.item;
			        var selectFunc = args.selectFunc;
			        var validFunc = args.validFunc;
			        var isValid = true;
			        if (validFunc != null) {
			            isValid = validFunc(item);
			        }
			        if (isValid) {
			            item.count--;
			            if (item.count <= 0) {
			                for (var i = 0; i < fromList.length; ++i) {
			                    if (fromList[i] == item) {
			                        fromList.splice(i, 1);
			                        break;
			                    }
			                }
			            }
			            var cartInst = null;
			            for (var i = 0; i < toList.length; ++i) {
			                if (toList[i].id == item.id) {
			                    var isSelect = true;
			                    if (selectFunc != null) {
			                        isSelect = selectFunc(toList[i]);
			                    }
			                    if (isSelect) {
			                        cartInst = toList[i];
			                        break;
			                    }
			                }
			            }
			            if (cartInst == null) {
			                cartInst = {
			                    id: item.id,
			                    count: 0
			                };
			                toList.push(cartInst);
			            }
			            cartInst.count++;
			        }
			    }
			    $scope.$on("buy", function(e, item) {
			        transaction({
			            from: $scope.itemInstList,
			            to: $scope.cartInstList,
			            item: item._inst,
			            selectFunc: function(it) {
			                return it.count < 3;
			            },
			            validFunc: function(it) {
			                var itemInfo = $scope.getItemInfo(it.id);
			                if (itemInfo.cost <= $scope.money) {
			                    $scope.money -= itemInfo.cost;
			                    return true;
			                } else {
			                    return false;
			                }
			            }
			        });
			        $scope.updateItemList();
			        $scope.updateCartList();
			    });
			    $scope.$on("sell", function(e, item) {
			        transaction({
			            from: $scope.cartInstList,
			            to: $scope.itemInstList,
			            item: item._inst,
			            selectFunc: function(it) {
			                return it.count < 2;
			            },
			            validFunc: function(it) {
			                var itemInfo = $scope.getItemInfo(it.id);
			                $scope.money += Math.round(itemInfo.cost * 0.5);
			                return true;
			            }
			        });
			        $scope.updateItemList();
			        $scope.updateCartList();
			    });
			}]);
		</script>
	</body>
</html>
```
1) script at element in plate will attach ng-controller and will be replaced to ng-eventName="functionName(arguments)"<br>
2) css in plate will affect only current plates DOM<br>
3) css has inhert="true" attribute in plate will affect from self to all child DOM<br>
4) plates in plate will be replaced(recursive problem is not solved yet)<br>
5) plates in DOM can have attribute & child, and they will be inserted & evaludated<br>
<br>
# /revel : revel intergration
add this code into "app/init.go" file<br>
<br>
imoprt _ "github.com/blackss2/goplate"<br>
<br>
"app/goplates" : file for plates(excpet views folder)<br>
"app/goplates/views" : files -> compile -> create result file at "app/views"<br>
<br>
exmaple : https://github.com/blackss2/goplate_revel_example<br>
