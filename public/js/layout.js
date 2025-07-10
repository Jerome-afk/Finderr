document.addEventListener("DOMContentLoaded", function() {
    // Toggle the sidebar
    const sidebar = document.querySelector(".sidebar");
    if (sidebar) {
        // Keep it minimised by default
        sidebar.classList.add("minimised");

        // Toggle on click
        sidebar.addEventListener("click", function() {
            this.classList.toggle("minimised");
        });
    }
});