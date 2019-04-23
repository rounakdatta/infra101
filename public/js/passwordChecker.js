function check(arg){
    var id = arg.getAttribute('id');
    var value = arg.value;
    
   if (id == 'email') {

       if ((/@srmuniv.edu.in\s*$/.test(value))) {
           document.getElementById('pwdMsg').style.color = 'green';
           document.getElementById('pwdMsg').innerHTML = 'great!';		
       } else {
           document.getElementById('pwdMsg').style.color = 'red';
           document.getElementById('pwdMsg').innerHTML = 'only @srmuniv.edu.in emails please!';			
   }

   }

   if (id == 'pwd') {
       if (value.length >= 6) {
       document.getElementById('pwdMsg').style.color = 'green';
       document.getElementById('pwdMsg').innerHTML = 'great!';
       
    } else {
       document.getElementById('pwdMsg').style.color = 'red';
       document.getElementById('pwdMsg').innerHTML = 'not strong enough';	 
    }

   }


    if ((/@srmuniv.edu.in\s*$/.test(document.getElementById('email').value)) && document.getElementById('pwd').value.length >= 6) {
        document.getElementById('submitButton').disabled = false;
    } else {
       document.getElementById('submitButton').disabled = true;
    }

   } 
