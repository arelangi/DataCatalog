$(document).ready(function() {
    $.get('/api', function(data, textStatus, jqXHR) {
        alert('status' + textStatus + ', data: ' + data);
    })
});