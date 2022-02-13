// Promt is the Js module for alerts, nofifications, and custom popup
function Prompt() {
	let toast = function (c) {
		const { msg = "", icon = "success" } = c;

		const Toast = Swal.mixin({
			toast: true,
			title: msg,
			position: "top-end",
			icon: icon,
			showConfirmButton: false,
			timer: 3000,
			timerProgressBar: true,
			didOpen: (toast) => {
				toast.addEventListener("mouseenter", Swal.stopTimer);
				toast.addEventListener("mouseleave", Swal.resumeTimer);
			},
		});

		Toast.fire({});
	};

	let success = function (c) {
		const { msg = "", title = "", footer = "" } = c;
		Swal.fire({
			icon: "success",
			title: title,
			text: msg,
			footer: footer,
		});
	};

	let error = function (c) {
		const { msg = "", title = "", footer = "" } = c;
		Swal.fire({
			icon: "error",
			title: title,
			text: msg,
			footer: footer,
		});
	};

	async function custom(c) {
		//https://stackoverflow.com/questions/54662735/what-does-the-statement-const-tz-msg-this-state-mean-in-the-below-code
		const { msg = "", title = "", icon = "", showConfirmButton = true } = c;

		// Take all data from return array (value method) of Swal.fire and assign to result
		const { value: result } = await Swal.fire({
			icon: icon,
			title: title,
			html: msg,
			backdrop: false,
			focusConfirm: false,
			showCancelButton: true,
			showConfirmButton: showConfirmButton,
			preConfirm: () => {
				return [
					document.getElementById("start").value,
					document.getElementById("end").value,
				];
			},

			willOpen: () => {
				if (c.willOpen !== undefined) {
					c.willOpen();
				}
			},
			// ASync runs after popup has been shown up in the screen
			didOpen: () => {
				if (c.didOpen !== undefined) {
					c.didOpen();
				}
			},
		});

		// Processing Swal data after client submit the date for a room
		if (result) {
			// if the client dont hit the cancel button, check to see if we have any actual values
			if (result.dismiss !== Swal.DismissReason.cancel) {
				if (result.value !== "") {
					if (c.callback !== undefined) {
						c.callback(result);
					}
				} else {
					c.callback(false);
				}
			} else {
				c.callback(false);
			}
		}
	}

	return {
		toast: toast,
		success: success,
		error: error,
		custom: custom,
	};
}
