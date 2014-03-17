
function showClock(div)
{
    var today=new Date();
    var h=today.getHours();
    var m=today.getMinutes();
    // add a zero in front of numbers<10
    if (m<10){
	m="0"+m
    }
    div.innerHTML=h+":"+m;
    t=setTimeout(function(){
	showClock(div)
    },30000);
}





