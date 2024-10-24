window.addEventListener("load", function () {
  let menuButton = document.getElementById("menu-button");
  let menu = document.getElementById("menu");
  let fakeBody = document.getElementById("body");
  let footer = document.getElementById("footer");

  fakeBody.addEventListener("click", () => {
    menu.classList.add("hidden");
    menu.classList.remove("transform-none");
    fakeBody.classList.remove("blur-lg");
    footer.classList.remove("blur-lg");
  });

  menuButton.addEventListener("click", () => {
    let timer = null;
    if (timer !== null) {
      clearTimeout(timer);
    }
    timer = setTimeout(function () {
      menu.classList.remove("hidden");
      menu.classList.add("transform-none");
      fakeBody.classList.add("blur-lg");
      footer.classList.add("blur-lg");
    }, 25);
  });
});
