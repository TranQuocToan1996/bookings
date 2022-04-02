function startSweetAlertForRoomBooking(room_id, csrf_token) {
	document
		.getElementById("check-availability-button")
		.addEventListener("click", function () {
			let html = `
                    <form id="check-availability-form" action="" method="post" novalidate class="needs-validation"
                        <div class="row">
                            <div class="col">
                                <div class="row" id="reservation-dates-modal">
                                    <div class="col">
                                        <input required disabled class="form-control" type="text" name="start" id="start" placeholder="Start date" autocomplete="off">
                                    </div>
                                    <div class="col">
                                        <input required disabled class="form-control" type="text" name="end" id="end" placeholder="End date" autocomplete="off">
                                    </div>
                                </div>
                            </div>
                        </div>
                    </form>
                    <br><br><br><br><br><br><br><br><br><br><br><br>
                    `;

			attention.custom({
				msg: html,
				title: "Choose your date",
				// Sync runs before datePicker has been popup shown up in the screen
				willOpen: () => {
					const elem = document.getElementById("reservation-dates-modal");
					const rangePicker = new DateRangePicker(elem, {
						format: "yyyy-mm-dd",
						showOnFocus: true,
						minDate: new Date(),
					});
				},
				// ASync runs after datePicker popup has been shown up in the screen
				didOpen: () => {
					document.getElementById("start").removeAttribute("disabled");
					document.getElementById("end").removeAttribute("disabled");
				},

				callback: function (result) {
					console.log(result);

					// Extract input tags (#start, #end) into formdata and append CSRF token before sending post request
					const form = document.getElementById("check-availability-form");
					let formData = new FormData(form);
					formData.append("csrf_token", csrf_token);
					// 1 is the room_id of general's Quarters
					formData.append("room_id", room_id);

					// Call handler
					fetch("/search-availability-json", {
						method: "post",
						body: formData,
					})
						.then((response) => response.json())
						.then((json) => {
							if (json.ok) {
								attention.custom({
									icon: "success",
									showConfirmButton: false,
									msg: `<strong>Room is available</strong>
                                        <div> <br> 
                            <a href="/book-room?id=${json.room_id}&s=${json.start_date}&e=${json.end_date}" 
                                        class="btn btn-primary">Book now</a>
                                        </div>`,
								});
							} else {
								attention.error({
									msg: `<strong>Room isn't availability</strong>`,
								});
							}
						});
				},
			});
		});
}