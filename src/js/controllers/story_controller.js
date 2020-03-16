import {
  Controller
} from 'stimulus';

export default class extends Controller {
  static targets = ['points', 'voterWrapper', 'voter', 'saver'];

  upvote(event) {
    const model = {
      UserID: parseInt(this.data.get('userid')),
      StoryID: parseInt(this.data.get('storyid')),
      VoteType: 1, // upvote
    }
    this.sendRequest(event, '/stories/vote', model, (res) => {
      if (res.Result === 'Voted') {
        this.voterTarget.setAttribute('data-action', 'click->story#removeUpvote');
        this.voterWrapperTarget.classList.add('upvoted');
        const currentPoints = this.data.get('points');
        const newPoints = parseInt(currentPoints) + 1;
        this.data.set('points', newPoints);
        this.pointsTarget.innerHTML = `${newPoints} points by `;
      }
    })
  }

  removeUpvote(event) {
    const model = {
      UserID: parseInt(this.data.get('userid')),
      StoryID: parseInt(this.data.get('storyid')),
      VoteType: 1, // upvote
    }
    this.sendRequest(event, '/stories/remove/vote', model, (res) => {
      if (res.Result === 'Unvoted') {
        this.voterTarget.setAttribute('data-action', 'click->story#upvote');
        this.voterWrapperTarget.classList.remove('upvoted');
        const currentPoints = this.data.get('points');
        const newPoints = parseInt(currentPoints) - 1;
        this.data.set('points', newPoints);
        this.pointsTarget.innerHTML = `${newPoints} points by `;
      }
    })
  }

  downvote(event) {

  }

  removeDownvote(event) {

  }

  save(event) {
    const model = {
      UserID: parseInt(this.data.get('userid')),
      StoryID: parseInt(this.data.get('storyid')),
    }
    this.sendRequest(event, '/stories/save', model, (res) => {
      if (res.Result === 'Saved') {
        this.saverTarget.setAttribute('data-action', 'click->story#unsave');
        this.saverTarget.innerHTML = 'unsave';
      }
    })
  }

  unsave(event) {
    const model = {
      UserID: parseInt(this.data.get('userid')),
      StoryID: parseInt(this.data.get('storyid')),
    }
    this.sendRequest(event, '/stories/unsave', model, (res) => {
      if (res.Result === 'Unsaved') {
        this.saverTarget.setAttribute('data-action', 'click->story#save');
        this.saverTarget.innerHTML = 'save';
      }
    })
  }

  sendRequest(event, url, model, onSuccess) {
    event.preventDefault();
    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }
    fetch(url, {
        method: 'POST',
        body: JSON.stringify(model)
      })
      .then(res => {
        return res.json();
      })
      .then(res => {
        onSuccess(res);
      });
  }
}