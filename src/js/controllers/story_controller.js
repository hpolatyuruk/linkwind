import {
  Controller
} from 'stimulus';

export default class extends Controller {
  static targets = ['points', 'voterWrapper', 'voter', 'saver'];

  upvote(event) {
    this.sendRequest(event, '/stories/upvote', (data) => {
      if (data.Result === 'Upvoted') {
        this.voterTarget.setAttribute('data-action', 'click->story#unvote');
        this.voterWrapperTarget.classList.add('upvoted');
        const currentPoints = this.data.get('points');
        const newPoints = parseInt(currentPoints) + 1;
        this.data.set('points', newPoints);
        this.pointsTarget.innerHTML = `${newPoints} points by `;
      }
    })
  }

  unvote(event) {
    this.sendRequest(event, '/stories/unvote', (data) => {
      if (data.Result === 'Unvoted') {
        this.voterTarget.setAttribute('data-action', 'click->story#upvote');
        this.voterWrapperTarget.classList.remove('upvoted');
        const currentPoints = this.data.get('points');
        const newPoints = parseInt(currentPoints) - 1;
        this.data.set('points', newPoints);
        this.pointsTarget.innerHTML = `${newPoints} points by `;
      }
    })
  }

  save(event) {
    this.sendRequest(event, '/stories/save', (data) => {
      if (data.Result === 'Saved') {
        this.saverTarget.setAttribute('data-action', 'click->story#unsave');
        this.saverTarget.innerHTML = 'unsave';
      }
    })
  }

  unsave(event) {
    this.sendRequest(event, '/stories/unsave', (data) => {
      if (data.Result === 'Unsaved') {
        this.saverTarget.setAttribute('data-action', 'click->story#save');
        this.saverTarget.innerHTML = 'save';
      }
    })
  }

  sendRequest(event, url, onSuccess) {
    event.preventDefault();
    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }
    const userID = this.data.get('userid');
    const storyID = this.data.get('storyid');
    fetch(url, {
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
        onSuccess(data);
      });
  }
}