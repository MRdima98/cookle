window.addEventListener("load", function () {
  let menuButton = document.getElementById("menu-button");
  let menu = document.getElementById("menu");
  let fakeBody = document.getElementById("body");

  fakeBody.addEventListener("click", () => {
    menu.classList.remove("transform-none");
    fakeBody.classList.remove("blur-lg");
  });

  menuButton.addEventListener("click", () => {
    let timer = null;
    if (timer !== null) {
      clearTimeout(timer);
    }
    timer = setTimeout(function () {
      menu.classList.add("transform-none");
      fakeBody.classList.add("blur-lg");
    }, 25);
  });
});
