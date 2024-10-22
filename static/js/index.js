window.addEventListener("load", function () {
  let menuButton = document.getElementById("menu-button");
  let menu = document.getElementById("menu");

  menuButton.addEventListener("click", () => {
    menu.style.visibility = "visible";
  });

  addEventListener("visibilitychange", () => {
    console.log("spooky");
  });
});
