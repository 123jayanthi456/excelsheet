// $(document).ready(function() {
// 	// Retrieve data from the server
// 	$.get("/persons", function(data) {
// 		$.each(data, function(i, person) {
// 			var row = $("<tr>").appendTo($("#personsTable tbody"));
// 			$("<td>").text(person.Name).appendTo(row);
// 			$("<td>").text(person.Phnno).appendTo(row);
// 			$("<td>").text(person.Email).appendTo(row);
// 			$("<td>").text(person.City).appendTo(row);
// 			$("<td>").text(person.State).appendTo(row);
// 			var deleteButton = $("<button>").text("Delete").appendTo($("<td>").appendTo(row));
// 			deleteButton.click(function() {
// 				$.ajax({
// 					url: "/persons/" + person.Name,
// 					type: "DELETE",
// 					success: function() {
// 						row.remove();
// 					}
// 				});
// 			});
// 		});
// 	});

// 	// Upload data to the server
// 	$("form").submit(function(e) {
// 		e.preventDefault();
// 		var formData = new FormData($(this)[0]);
// 		$.ajax({
// 			url: "/upload",
// 			type: "POST",
// 			data: formData,
// 			processData: false,
// 			contentType: false,
// 			success: function() {
// 				location.reload();
// 			}
// 		});
// 	});
// });
