$(document).ready(function(){

var name;
var dob;
var coverageName;
var coverageType;
var carrierId;
var group;
var plan;
var subscriberId;
var primary;
var startDate;
var endDate;
var employer;
var dependents;
var premium;


$('#searchEmployee').on('click', function(){
  $.get('https://insHfcBluemixTest.mybluemix.net/getCoverages?subscriberID=ba2345', function(result){
		//console.log(result['Coverages'][0].subsciberName);
    //console.log(result.Coverages[0].subscriberDOB);
	 //console.log(result['Coverages'][0]['subscriberName']);
	 $('#name').val(result.SubscriberName);
	 $('#dob').val(result.SubscriberDOB);
   $('#subscriberId').val(result.SubscriberID);
   $('#coverageType').val(result.CoverageType);
   $('#coverageName').val(result.CoverageName);
   $('#startDate').val(result.StartDate);
   $('#endDate').val(result.EndDate);
   $('#plan').val(result.PlanCode);
   $('#dependents').val(result.Dependents.length);
   $('#premium').val(result.Premium);



    info=result;
      //console.log(result);
      console.log(info);

		 }, 'json');

});



});
