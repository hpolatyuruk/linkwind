import {
  Controller
} from 'stimulus';

export default class extends Controller {
  static targets = ['points', 'voterWrapper', 'voter'];

  upvote(event) {
    event.preventDefault();
    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }

    const userID = this.data.get('userid');
    const storyID = this.data.get('storyid');
    fetch('/stories/upvote', {
        method: 'POST',
        body: JSON.stringify({
          StoryID: parseInt(storyID),
          UserID: parseInt(userID)
        })
      })
      .then(res => {
        return res.json();
      })
      .then(data => {
        console.log(data);
        if (data.Result === 'Upvoted') {
          this.voterTarget.setAttribute('data-action', 'click->story#unvote');
          this.voterWrapperTarget.classList.add('upvoted');
          const currentPoints = this.data.get('points');
          const newPoints = parseInt(currentPoints) + 1;
          this.data.set('points', newPoints);
          this.pointsTarget.innerHTML = `${newPoints} points by `;
        }
      });
  }

  unvote() {
    event.preventDefault();
    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }

    const userID = this.data.get('userid');
    const storyID = this.data.get('storyid');
    fetch('/stories/unvote', {
        method: 'POST',
        body: JSON.stringify({
          StoryID: parseInt(storyID),
          UserID: parseInt(userID)
        })
      })
      .then(res => {
        return res.json();
      })
      .then(data => {
        console.log(data);
        if (data.Result === 'Unvoted') {
          this.voterTarget.setAttribute('data-action', 'click->story#upvote');
          this.voterWrapperTarget.classList.remove('upvoted');
          const currentPoints = this.data.get('points');
          const newPoints = parseInt(currentPoints) - 1;
          this.data.set('points', newPoints);
          this.pointsTarget.innerHTML = `${newPoints} points by `;
        }
      });
  }
}