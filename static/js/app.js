// prompt handles all the SweetAlert
function prompt() {
    let toast = function(c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;

        Swal.fire({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })
    }

    let success = function(c) {
        const {
            title = "",
            text = "",
        } = c;

        Swal.fire({
            icon: "success",
            title: title,
            text: text
            })
    }

    let error = function(c) {
        const {
            title = "",
            text = "",
        } = c;

        Swal.fire({
            icon: "error",
            title: title,
            text: text
            })
    }

    async function custom(c) {
        const {
            icon = "", // Icon that will be displayed
            text = "", // HTML text that will be displayed
            title = "", // Title that will be displayed
            showConfirmButton = true,
        } = c;

        const { value: formValues } = await Swal.fire({
            icon: icon,
            title: title,
            html: text,
            backdrop: true,
            focusConfirm: true,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                document.getElementById('arrivalDate_2').value,
                document.getElementById('departureDate_2').value
                ]
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen()
                }
            }
        })

        // Checks if there is soe result comming back from the dialog box
        if (formValues) {
            // If the result is NOT the cancel button
            if (formValues.dismiss !== Swal.DismissReason.cancel) {
                // If the result is not empty
                if (formValues.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(formValues);
                    }
                } else {
                    c.callback(false)
                }
            } else {
                c.callback(false)
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}