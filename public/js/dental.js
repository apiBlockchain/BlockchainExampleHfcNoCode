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
var annualDeductableBal;
var annualBenefitMaxBalance;
//var result = {}


$('#search').on('click', function(){
  $('#popup').fadeIn();
  	window.subscriberId = $('#subscriberID').val();
    window.memberId = $('#memberID').val();
    console.log(window.subscriberId);
    console.log(window.memberId);
  $.get('https://insHfcBluemixTest.mybluemix.net/verifyEmployment?subscriberId='+window.subscriberId, function(result){
    if (result.Status == "Passed"){
      console.log("status passed");

        $('.employer-icon').addClass('employer-pass').removeClass('employee-check');
        $('.message-details1').html(result.Message).delay(5000).fadeOut('fast');
        secondCall();

          }
          else
          {
              console.log("no");
          }

		 }, 'json');

     function secondCall(){
     $.get('https://insHfcBluemixTest.mybluemix.net/verifyCoverage?subscriberId='+window.subscriberId+ '&memberId=' +window.memberId, function(result){


          $('.coverage-icon').delay(5500).queue(function(){
              $(this).addClass('dental-pass').dequeue();
              $(this).removeClass('dental-check');

            });

          $('.message-details2').hide(0).delay(5500).fadeIn().html(result.Message);


           if (result.Status == "Passed"){
             console.log("status passed");

              $('.check-icon').delay(11000).queue(function(){
                 $(this).addClass('verified-pass').dequeue();
                 $(this).removeClass('verified-check');

               });


               populate();
               $('#patientInfo').fadeIn();
                 }
                 else
                 {
                     console.log("no");

                     $('.check-icon').delay(11000).queue(function(){
                        $(this).addClass('verified-check-no').dequeue();
                        $(this).removeClass('verified-check');

                      });


                 }

                 $



   		 }, 'json');
     }

     $('#close').click(function(){
       $('#popup').hide();


     })

       function populate(){

         $.get('https://insHfcBluemixTest.mybluemix.net/getCoverages?subscriberID=ba2345', function(result){
         console.log(result);
         window.dentalInfo = result;
         //console.log(result.Coverages[0].subscriberDOB);
         //console.log(result['Coverages'][0]['subscriberName']);
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
         $('#annualBenefitMaximumBalance').val(result.AnnualBenefitMaximumBal);
         $('.deductible').html(result.AnnualDeductible);
         $('.annual-benefits').html(result.AnnualBenefitMaximum);
         $('.lifetime-max').html(result.LifetimeBenefitMaximum);
         $('.diagnostic').html(result.PreventiveCare);
         $('.minor-restorative').html(result.MinorRestorativeCare);
         $('.major-restorative').html(result.MajorRestorativeCare);
         $('.ortho-treatment').html(result.OrthodonticTreatment);
         $('.ortho-benefits').html(result.OrthodonticLifetimeBenefitMaximum);
         name = result.subscriberName;
         dob = result.subscriberDOB;
         employer = result.EmployerID;
         employee = result.subscriberID;

         info=result;

         console.log(info);


         }, 'json');

         //});



       }

});








});
