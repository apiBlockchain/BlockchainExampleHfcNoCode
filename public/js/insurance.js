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


$('#search').on('click', function(){
  $.get('https://insHfcBluemixTest.mybluemix.net/getCoverages?subscriberID=ba2345', function(result){
		console.log(result);
    window.dentalInfo = result;

   $('#name').val(result.SubscriberName);
   $('#dob').val(result.SubscriberDOB);
   $('#subscriberId').val(result.SubscriberID);
   $('#primary').val(result.IsPrimary);
   $('#coverageType').val(result.CoverageType);
   $('#coverageName').val(result.CoverageName);
   $('#carrierId').val(result.CarrierID);
   $('#startDate').val(result.StartDate);
   $('#endDate').val(result.EndDate);
   $('#group').val(result.GroupNum);
   $('#plan').val(result.PlanCode);
   $('#employer').val(result.EmployerID);
   $('#dependents').val(result.Dependents.length);
   $('#premium').val(result.Premium);
   $('#annualDeductableBalance').val(result.AnnualDeductibleBal);
   $('#annualBenefitMaxBalance').val(result.AnnualBenefitMaximumBal);
   $('.deductible').html(result.AnnualDeductible);
   $('.annual-benefits').html(result.AnnualBenefitMaximum);
   $('.lifetime-max').html(result.LifetimeBenefitMaximum);
   $('.diagnostic').html(result.PreventiveCare);//showing up weird in the object
   $('.minor-restorative').html(result.MinorRestorativeCare);
   $('.major-restorative').html(result.MajorRestorativeCare);
   $('.ortho-treatment').html(result.OrthodonticTreatment);
   $('.ortho-benefits').html(result.OrthodonticLifetimeBenefitMaximum);
	name = result.subscriberName;
    dob = result.subscriberDOB;
    employer = result.EmployerID;
    employee = result.subscriberID;
    // var date = new Date(data.DOB);
    // var formattedDate = date.getFullYear() + "-" +leadingZero(date.getMonth()+1)+ "-" + date.getDate();

    info=result;
      //console.log(result);
      console.log(info);


		 }, 'json');

});



$('#updateInfo').click(function(){

  var subsriber = $('#subscriberId').val();
  var annualDeductableBal = $('#annualDeductableBalance').val();
  var annualBenefitMaxBal = $('#annualBenefitMaxBalance').val();

  dentalInfo['AnnualDeductibleBal'] = annualDeductableBal;
  dentalInfo['AnnualBenefitMaximumBal'] = annualBenefitMaxBal;





  	$.ajax({
    	url: " https://insHfcBluemixTest.mybluemix.net/updateCoverage",
    	type: "POST",
    	data:  JSON.stringify(dentalInfo),
    	contentType: "application/json",
    	success: function(response){

        //console.log(annualDeductableBal, "this is deductable");
        //console.log(annualBenefitMaxBal, "this is bal");

      	console.log(response, "response");
      	console.log(dentalInfo['AnnualDeductibleBal']);
        console.log(dentalInfo['AnnualBenefitMaximumBal']);
   	 }
  	});

});



});
