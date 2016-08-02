var upel = $('.hero.is-primary')
var downel = $('.hero.is-danger')
var time = 1250

function showstatus(status) {
	if (status == "up") {
		upel.css("display","block")
		setTimeout(function(){
			upel.css("display","none")
		},time)
	} else if (status == "down"){
		downel.css("display","block")
		setTimeout(function(){
			downel.css("display","none")
		},time)
	} else{
		upel.css("display","none")
		downel.css("display","none")
	}
}

function sitestatus(url){
	var list
	var statuses = ""
	$.getJSON('/list', function(data){
	  for (var i = 0;i<data.length;i++) {
	  	var url = data[i]
		$.post('/isdown',{"url":url},function(res){
			if (res == "false"){
				status = "up"
			}else {
				status = "down"
			}
		})
	}
	})
	return status
}

$("#isdown-btn").on('click',function(e){
	var site = $("#url").val()
	var status = sitestatus(site)
	showstatus(status)
});