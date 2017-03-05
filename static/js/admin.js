$(document).ready(function() {
  setInterval(function() {
    update_page();
  }, 2000);
});

function update_page() {
  $("#main").load("/ajax/adminplaylist");
}