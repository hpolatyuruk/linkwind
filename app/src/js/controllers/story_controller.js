import { Controller } from 'stimulus';

export default class extends Controller {
  static targets = [];

  upvote(event) {
    event.preventDefault();
    let isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }

    let userID = this.data.get('userid');
    let storyID = this.data.get('storyid');
    fetch('/stories/upvote', {
      method: 'POST',
      body: JSON.stringify({
        StoryID: parseInt(storyID),
        UserID: parseInt(userID)
      })
    }).then(res => {
      if (res.status != 200) {
        console.log(res);
      }
    });
  }

  unvote() {}
}
