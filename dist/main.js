// Simple form validation
document.addEventListener('DOMContentLoaded', function() {
    const forms = document.querySelectorAll('.contact-form');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            let isValid = true;
            const requiredFields = form.querySelectorAll('[required]');

            requiredFields.forEach(field => {
                if (!field.value.trim()) {
                    isValid = false;
                    // You might want to add a class to highlight the invalid field
                    field.style.borderColor = 'red';
                } else {
                    field.style.borderColor = ''; // Reset border color
                }
            });

            if (!isValid) {
                e.preventDefault(); // Stop form submission
                // You could also display a general error message to the user
                console.log('Please fill out all required fields.');
            }
        });
    });
});
