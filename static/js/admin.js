var domain;

$(document).ready(function() {
    domain = window.location.href.substring(0, window.location.href.lastIndexOf("/"));
    admin_playlist_refresh();
});

function admin_playlist_refresh() {
    $("#main").load(domain + "/ajax/admin/playlist", function(response, status) {
        // Set timeout dependant on success or error of previous request
        const timeout = status === "success" ? 2000 : 15000;
        window.setTimeout(function() {
            admin_playlist_refresh();
        }, timeout);
    });
}
