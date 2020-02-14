import {
  Controller
} from 'stimulus';

export default class extends Controller {
  static targets = ['replyForm', 'replyText', 'replyOutput', 'voter', 'voterWrapper', 'points'];

  showReplyBox(event) {
    event.preventDefault();
    console.log('clicked');
    const replyForm = this.replyFormTarget;
    replyForm.classList.add('flex-row');
    replyForm.classList.add('w-full');
    replyForm.classList.add('ml-10');
    replyForm.classList.add('mt-2');
    replyForm.innerHTML =
      "<div class='flex-row w-full'><textarea data-target='comment.replyText' id='reply' name='reply' rows='5' columns='25' class='bg-gray-200 appearance-none border-2 border-gray-200 rounded w-1/2 py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-purple-500 text-sm'></textarea></div>";
    replyForm.innerHTML += "<div class='flex-row w-full'>";
    replyForm.innerHTML +=
      "<button data-action='click->comment#reply' class='shadow bg-gray-200 hover:bg-gray-400 focus:shadow-outline focus:outline-none text-gray-600 text-sm font-semibold py-1 px-5 rounded text-sm' type='button'>Post</button>";
    replyForm.innerHTML +=
      " <button data-action='click->comment#cancel' class='shadow bg-gray-200 hover:bg-gray-400 focus:shadow-outline focus:outline-none text-gray-600 text-sm font-semibold py-1 px-5 rounded text-sm' type='button'>Cancel</button>";
    replyForm.innerHTML += '</div>';
  }

  reply(event) {
    event.preventDefault();
    const replyText = this.replyTextTarget.value;

    if (replyText === '') {
      return;
    }

    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }

    const storyID = this.data.get('storyid');
    const userName = this.data.get('username');
    const parentCommentID = this.data.get('commentid');

    fetch('/comments/reply', {
        method: 'POST',
        body: JSON.stringify({
          ParentCommentID: parseInt(parentCommentID),
          StoryID: parseInt(storyID),
          ReplyText: replyText
        })
      })
      .then(res => {
        if (res.ok) {
          return res.text();
        }
        return ''
      })
      .then(res => {
        if (res != '') {
          this.removeReplyForm()
          this.replyOutputTarget.innerHTML = res
        }
      });
  }

  cancel(event) {
    event.preventDefault();
    this.removeReplyForm()
  }

  upvote(event) {
    event.preventDefault();
    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }

    const userID = this.data.get('userid');
    const commentID = this.data.get('commentid');
    fetch('/comments/upvote', {
        method: 'POST',
        body: JSON.stringify({
          UserID: parseInt(userID),
          CommentID: parseInt(commentID)
        })
      })
      .then(res => {
        return res.json();
      })
      .then(data => {
        console.log(data);
        if (data.Result === 'Upvoted') {
          this.voterTarget.setAttribute('data-action', 'click->comment#unvote');
          this.voterWrapperTarget.classList.add('upvoted');
          const currentPoints = this.data.get('points');
          const newPoints = parseInt(currentPoints) + 1;
          this.data.set('points', newPoints);
          this.pointsTarget.innerHTML = ` | ${newPoints} points`;
        }
      });
  }

  unvote(event) {
    event.preventDefault();
    const isAuthenticated = this.data.get('isauthenticated') == 'true';
    if (isAuthenticated === false) {
      window.location = '/signin';
      return;
    }

    const userID = this.data.get('userid');
    const commentID = this.data.get('commentid');
    fetch('/comments/unvote', {
        method: 'POST',
        body: JSON.stringify({
          UserID: parseInt(userID),
          CommentID: parseInt(commentID)
        })
      })
      .then(res => {
        return res.json();
      })
      .then(data => {
        console.log(data);
        if (data.Result === 'Unvoted') {
          this.voterTarget.setAttribute('data-action', 'click->comment#upvote');
          this.voterWrapperTarget.classList.remove('upvoted');
          const currentPoints = this.data.get('points');
          const newPoints = parseInt(currentPoints) - 1;
          this.data.set('points', newPoints);
          this.pointsTarget.innerHTML = ` | ${newPoints} points`;
        }
      });
  }

  removeReplyForm = () => {
    const replyForm = this.replyFormTarget;
    replyForm.innerHTML = '';
    replyForm.className = '';
  }
}