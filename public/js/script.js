
$(document).ready(function(){



  $('.logo').click(function(){
     window.location.href='index.html';
  })

  $('.enroll-btn-link').click(function(){
       $('#popup').show();
  })

  $('.verify').click(function(){
    $('#popup').hide();
    window.location.href='enroll.html';
  })
  $('#selectPlan').click(function(){
    $('.confirm').show();
  })
  $('.confirm-btn').click(function(){
    $('.confirm').hide();
    window.location.href='policies.html';

  })
  $('.blue-link').click(function(){
    $('.long-desc').show();
    $('.short-desc').hide();

  })

  $('.blue-link-expand').click(function(){

    $('.long-desc').hide();
      $('.short-desc').show();
  })
  $('.blue-link2').click(function(){
    $('.long-desc2').show();
    $('.short-desc2').hide();

  })

  $('.blue-link-expand2').click(function(){

    $('.long-desc2').hide();
      $('.short-desc2').show();
  })
  $('.blue-link3').click(function(){
    $('.long-desc3').show();
    $('.short-desc3').hide();

  })

  $('.blue-link-expand3').click(function(){

    $('.long-desc3').hide();
      $('.short-desc3').show();
  })


  function leadingZero(n){
    		if(n < 10){
    			n = '0'+n;
    		}
    		return n;
    	}

  //var name;
  var email;
  var info;
  var street;
  var city;
  var state;
  var zip;
  var phone;
  var dob;
  var newInfo;

  var employer;
  var employee;
  var coverageName;
  var coverageType;
  var premium;
  var dependents;







	$.get('https://insHfcBluemixTest.mybluemix.net/getEmployeeRecord?employeeId=294048', function(data){
		window.name = data.Name;
    street = data.Address.Street;
    city = data.Address.City;
    state = data.Address.State;
    zip = data.Address.Zip;
    phone = data.Phone;
		email = data.Email;
    dob = data.DOB;
    window.employer = data.EmployerID;
    window.employee = data.EmployeeID;
    var date = new Date(data.DOB);
    var formattedDate = date.getFullYear() + "-" +leadingZero(date.getMonth()+1)+ "-" + date.getDate();
    console.log(formattedDate);
    console.log(window.employer);
    info=data;
	   populate();


     newInfo = {
        "CoverageName":"",
        "CoverageType": "",
        "CarrierID": "AE340",
        "GroupNum": "G10023",
        "PlanCode": "DP001",
        "SubscriberID": "ba2345",
        "SubscriberName":window.name,
        "SubscriberDOB": formattedDate,
        "IsPrimary": "YES",
        "StartDate": "2017-01-01",
        "EndDate": "2017-12-31",
        "AnnualDeductible":"850",
        "AnnualBenefitMaximum":"5000",
        "LifetimeBenefitMaximum": "Lifetime Benefit Maximum - None",
        "PreventiveCare": "Diagnostic/Preventive Care - In network no charge, Out-of-network-20%",
        "MinorRestorativeCare": "Minor Restorative Care - In network 20% of negotiated fee, Out-of-network-20% of U&P",
        "MajorRestorativeCare": "Major Restorative Care - Not Covered",
        "OrthodonticTreatment": "Orthodontic Treatment - Not Covered",
        "OrthodonticLifetimeBenefitMaximum": "Orthodontic Lifetime Benefit Maximum - N/A",
        "AnnualDeductibleBal": "0",
        "AnnualBenefitMaximumBal":"0",
        "EmployeeID": "294048",
        "MemberID": "M-01",
        "EmployerID": "Global Industries",
        "Dependents": [
          {
            "MemberName": "Megan Sheen",
            "MemberID": "M-03",
            "MemberDOB": "08/20/1990",
            "SubscriberID": "ba2345"
          },
          {
            "MemberName": "Wade Sheen",
            "MemberID": "M-02",
            "MemberDOB": "08/20/1961",
            "SubscriberID": "ba2345"
          }
        ],
        "Premium": "250"
      }






		}, 'json');

	function populate(){


    console.log(info);


    $('.name').html(name);
    $('.dob').html(dob);
    $('.employer-id').html(window.employer);
    $('.employee-id').html(window.employee);
    $('.street').html(street);
    $('.city-state-zip').html(city + "," + state + "" +zip);
    $('.phone').html(phone);
    $('.email').html(email);

}








$('#selectPlan').on('click', function(){



  var type = $(this).data('plan-type');
  var planName = $(this).data('plan-name');
  var price = $(this).data('premium');

  newInfo['CoverageType'] = type;
  newInfo['CoverageName'] = planName;
  //newInfo['Premium'] = price;
  //newInfo['AnnualDeductible'] = "$850";
  //newInfo['AnnualBenefitMaximum'] ="$5000.00";

    console.log(newInfo);
    console.log(premium);
    

    $.ajax({
      url: " https://insHfcBluemixTest.mybluemix.net/addCoverage",
      type: "POST",
      data:  JSON.stringify(newInfo),
      contentType: "application/json",
      success: function(response){
        console.log(newInfo);
        console.log(response, "response");
        console.log(newInfo['CoverageType']);
        console.log(newInfo['SubscriberName']);
        console.log(newInfo['AnnualDeductible']);
        console.log(newInfo['Premium']);
        //console.log();


      }
    });
});









});
