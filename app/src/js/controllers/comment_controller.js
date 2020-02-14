import {
  Controller
} from 'stimulus';

export default class extends Controller {
  static targets = ['form', 'replyText', 'output'];

  showReplyBox(event) {
    event.preventDefault();
    console.log('clicked');
    const replyForm = this.formTarget;
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
          this.outputTarget.innerHTML = res
        }
      });
  }

  cancel(event) {
    event.preventDefault();
    this.removeReplyForm()
  }

  upvote(event) {
    event.preventDefault();
    console.log('upvoted');
  }

  unvote(event) {
    event.preventDefault();
    console.log('unvoted');
  }

  removeReplyForm = () => {
    const replyForm = this.formTarget;
    replyForm.innerHTML = '';
    replyForm.className = '';
  }
}