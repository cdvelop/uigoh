// Component: Navbar
(function() {
  const navbarShowBtn = document.querySelector('.navbar-show-btn');
  const navbarCollapseDiv = document.querySelector('.navbar-collapse');
  const navbarHideBtn = document.querySelector('.navbar-hide-btn');

  if (navbarShowBtn && navbarCollapseDiv && navbarHideBtn) {
    navbarShowBtn.addEventListener('click', function() {
      navbarCollapseDiv.classList.add('navbar-show');
    });

    navbarHideBtn.addEventListener('click', function() {
      navbarCollapseDiv.classList.remove('navbar-show');
    });
  }

  // Change search icon on window resize
  function changeSearchIcon() {
    const searchIcon = document.querySelector('.search-icon img');
    if (!searchIcon) return;

    let winSize = window.matchMedia("(min-width: 1200px)");
    if (winSize.matches) {
      searchIcon.src = "images/search-icon.png";
    } else {
      searchIcon.src = "images/search-icon-dark.png";
    }
  }

  window.addEventListener('resize', changeSearchIcon);
  changeSearchIcon();
})();
