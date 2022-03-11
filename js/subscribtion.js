let formations = null
let matters = null
let schools = null

osOnChangeTimerDelay = 3000;

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
  getInfos();
  
  // verify passwords
  
  document.getElementById('pwd').oninput = function() {
    comparePWD(document.getElementById('pwd').value, document.getElementById('pwd-confirm').value)
    
  };
  document.getElementById('pwd-confirm').oninput = function() {
    comparePWD(document.getElementById('pwd').value, document.getElementById('pwd-confirm').value)
  };
}



var rad = document.subscribe.account;
var prev = null;
for (var i = 0; i < rad.length; i++) {
    rad[i].addEventListener('change', function() {
        (prev) ? prev.value: null;
        if (this !== prev) {
            prev = this;
        }
        if (this.value == "student") {
          showFormation();
          showGit();
        } else {
          showMatter();
          removeGit();
        }                  
    });
}
function showFormation() {
  document.getElementById('studiesOrMatter').innerHTML = `
      <label for="formation">Formation</label>
      <select title="Vous devez choisir une formation" name="formation" id="formation" required>
      </select>
    `;
  options = `<option value="" selected disabled>--Choisissez une formation--</option>`;

  for (const formation of formations) {
    options = options+`<option value="`+formation+`">`+formation+`</option>`
  }
  document.getElementById('formation').innerHTML = options

}


function showMatter() {
  document.getElementById('studiesOrMatter').innerHTML = `
      <label for="matiere">Matière</label>
      <select title="Vous devez choisir une matière" name="matiere" id="matiere" required>
      </select>
  `;
  options = `<option value="" selected disabled>--Choisissez une matière--</option>`
  for (const matter of matters) {
    options = options+`<option value="`+matter+`">`+matter+`</option>`
  }
  document.getElementById('matiere').innerHTML = options
}


function showSchool() {
  document.getElementById('selectSchool').innerHTML = `
    <label for="campus">Campus</label>
    <select title="Vous devez choisir un campus" name="campus" id="campus" required>
      
      <option value="Paris 19">Paris 19</option>
      <option value="Melun">Melun</option>
      <option value="Rennes">Rennes</option>
    </select>
  `;
  options = `<option selected disabled value="">--Choisissez une école--</option>`
  for (const school of schools) {
    options = options+`<option value="`+school+`">`+school+`</option>`
  }
  document.getElementById('campus').innerHTML = options
}

function showGit() {
  document.getElementById('repoField').innerHTML = `
  <label for="repo"> Lien vers votre repository </label>
  <input title="Le nombre de caractères maximum est de 250" required maxlength="250" id="repo" name="repo" type="text" placeholder="ex: https://github.com/TavernierAlicia/">
  `;
}
function removeGit() {
  document.getElementById('repoField').innerHTML = ``;
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