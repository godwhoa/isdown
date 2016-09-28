if (!String.prototype.format) {
  String.prototype.format = function() {
    var args = arguments;
    return this.replace(/{(\d+)}/g, function(match, number) { 
      return typeof args[number] != 'undefined'
        ? args[number]
        : match
      ;
    });
  };
}

// you can pick which minion to test with
function ex_check_site(url,minion,testing_cb,down_cb,up_cb){
    testing_cb(minion)
    $.post("http://"+minion+ '/isdown', {
        url: url
    }, function(res,status) {
    	console.log(status,i)
        if (res == "false"){
        	up_cb(minion)
        	// console.log(minion,"site up!")
        }else{
        	down_cb(minion)
        	// console.log(minion,"down.")
        }
    })
}


// tcb: testing callback
// dcb: down callback
// ucb: up callback
function check_site(url,tcb,dcb,ucb) {
    $.getJSON('/list',
        function(list) {
        	console.log(list)
            for (var i = 0; i < list.length; i++) {
                var minion = list[i]
                tcb(minion)
            }
    });
}

var testing_template = '<div><img src="img/loader.svg"><a>{0} is testing...</a><br></div>'
var done_template = '<div><a>{0}</a><br></div>'

/* handles and displays state of minions. */
var MDisplay = function() {
	this.status = {}
	this.container = $(".results-container")
}

MDisplay.prototype.set = function(minion,status) {
	this.status[minion] = status
	this.draw()
}

MDisplay.prototype.draw = function() {
	this.container.empty()
	var self = this
	$.each(self.status, function(minion, status) {
		if (status == "testing"){
			self.container.append(testing_template.format(minion))
		}else if (status == "up"){
			self.container.append(done_template.format("Up!"))
		}else{
			self.container.append(done_template.format("Down."))
		}
	})
}

var mdisplay = new MDisplay()

$('.test-button').click(function(){
	var url = $(".input-container > input").val()
	// hackish way of clearning out data
	mdisplay = new MDisplay()
	check_site(url,function(minion){
		mdisplay.set(minion,"testing")
	},function(minion){
		mdisplay.set(minion,"down")
	},function(minion){
		mdisplay.set(minion,"up")
	})
})