# goplate
# / : template style Front-end preprocessor for angularjs
- Working In Progress<br>
1) embeded css<br>
css block only affect within current plate block(if inherit=true, subplate blocks are also affeted by current css block)<br>
<br>
2) embeded javascript(intergrated into angularjs)<br>
script block will be change to angularjs event, and parents of script block will be controller<br>
attribute inject of script will inject module into current Controller<br>
<br>
3) each plate turn to html element<br>
if you have plate has name subject, when you use subject element, it will be replaced with plate<br>
all other inline html, script plate element have will override original plate<br>

## Before Compile
```html
<plate name="subject">
	<css>
	.ABC {
		background-color: black;
	}
	</css>
	<li class="ABC">
		<span class="BCD" sid="{{subject.Sid}}">{{subject.Name}}<div EF="123">Subject Div Text</div></span>
		<script event="click(subject)" inject="$http">
		$scope.$emit('selectSubject', subject.Sid);
		//alert(subject.Sid);
		</script>
		<script event="mouseover(subject)" inject="$q">
		//alert(subject.Name);
		</script>
	</li>
</plate>

<!-- subjectList.html -->
<plate name="subjectList">
	<css inherit="true">
	.ABC {
		background-color: red;
		.BCD {
			background-color: orange;
		}
		&:hover {
			background-color: blue;
		}
	}
	</css>
	<ul>
		<subject ng-repeat="subject in subjectList" ng-if="subjectList.length > 0">
			<div>
				<label>Inline Label</label>
				<script event="mouseover(subject)" inject="$http">
				//alert(subject.Sid);
				</script>
			</div>
		</subject>
		<script>
		//Load Event
		$scope.subjectList = [{
			Name : "TestName",
			Sid : "123"
		}];
		$scope.$on('selectSubject', function(e, Sid) {
			alert(Sid);
		});
		</script>
	</ul>
</plate>

<subjectList>
</subjectList>
```
## After Compile
```HTML
<!-- subject.html -->
<html>
	<head>
	</head>
	<body ng-app="">
		<!-- subjectList.html -->
		<ul ng-controller="Ctrl_3">
			<li class="ABC  genclass_1_1 genclass_2_1" ng-controller="Ctrl_1" ng-click="EventHandler_1_1(subject)" ng-mouseover="EventHandler_1_2(subject)" ng-repeat="subject in subjectList" ng-if="subjectList.length &gt; 0">
				<span class="BCD  genclass_2_2" sid="{{subject.Sid}}">
					{{subject.Name}}
					<div ef="123">
						Subject Div Text
					</div>
				</span>
				<div ng-controller="Ctrl_2" ng-mouseover="EventHandler_2_1(subject)">
					<label>
						Inline Label
					</label>
				</div>
			</li>
		</ul>
	</body>
</html>
<script>
	function Ctrl_1($scope, $element, $http, $q) {
	    $scope.EventHandler_1_1 = function(subject) {
	        $scope.$emit('selectSubject', subject.Sid);
	        //alert(subject.Sid);
	    }
	    $scope.EventHandler_1_2 = function(subject) {
	        //alert(subject.Name);
	    }
	}
	function Ctrl_2($scope, $element, $http) {
	    $scope.EventHandler_2_1 = function(subject) {
	        //alert(subject.Sid);
	    }
	}
	function Ctrl_3($scope, $element) {
	    //Load Event
	    $scope.subjectList = [{
	        Name: "TestName",
	        Sid: "123"
	    }];
	    $scope.$on('selectSubject', function(e, Sid) {
	        alert(Sid);
	    });
	}
</script>
<style>
	.ABC.genclass_1_1 {
		background-color: black;
	}
	.ABC.genclass_2_1 {
		background-color: red;
	}
	.ABC.genclass_2_1 .BCD.genclass_2_2 {
		background-color: orange;
	}
	.ABC.genclass_2_1:hover {
		background-color: blue;
	}
</style>
```

# /revel : revel intergration
- Working In Progress<br>
goplate.Render(controller, v ...interface{}) revel.Result

