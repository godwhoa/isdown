function check_site(url) {
    $.getJSON('/list',
        function(list) {
            for (var i = 0; i < list.length; i++) {
                var minion = list[i]
                $.post(minion + '/isdown', {
                    url: 'payload'
                }, function(res) {
                    if (res == "false"){
                    }else{

                    }
                })
            }
        }
    });
}