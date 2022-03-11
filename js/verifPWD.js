let formations = null
let matters = null
let schools = null


const getInfos = _ => fetch('http://localhost:9999/subInfos').then(res => res.json()).then(res => showInfos(res));

const showInfos = (infos) => {
  formations = infos.Formations;
  matters = infos.Matters;
  schools = infos.Schools;
  showSchool();
  showFormation();
  showGit();
}

window.onload = function() {

  // verify passwords
  
  document.getElementById('pwd').oninput = function() {
    comparePWD(document.getElementById('pwd').value, document.getElementById('pwd-confirm').value)
    
  };
  document.getElementById('pwd-confirm').oninput = function() {
    comparePWD(document.getElementById('pwd').value, document.getElementById('pwd-confirm').value)
  };
}

function comparePWD(pwd, confirm) {
  if (pwd != confirm) {
    document.getElementById('infos-pwd').textContent = `Les mots de passe ne sont pas identiques`
    document.getElementById('pwd').setAttribute("pattern", "[A]{1000,1000}")
  } else {
    document.getElementById('infos-pwd').textContent = ``
    document.getElementById('pwd').removeAttribute("pattern")
  }
}