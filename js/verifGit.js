const regexGitHub = /github.com/i;

window.addEventListener("load", _ => {
    var repos = document.querySelectorAll('[class=repo]');
    // http request
    for (var i = 0; i < repos.length; i++) {
        let repo = repos[i]
        link = repo.textContent.replace(regexGitHub, 'api.github.com/repos')
        fetch(link).then(res => {
            if (res.ok) {
                console.log("green")
                repo.style.color='#40A133'
            } else {
                console.log("red")
                repo.style.color='#A53636'
            }
        }).catch(_ => {
            console.log("red")
            repo.style.color='#A53636'
        })
    }
})