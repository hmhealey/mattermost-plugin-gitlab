import {connect} from 'react-redux';

import {getUser} from 'mattermost-redux/selectors/entities/users';

import UserAttribute from './user_attribute.jsx';

function mapStateToProps(state, ownProps) {
    const id = ownProps.user ? ownProps.user.id : '';
    const user = getUser(state, id);

    let username;
    if (user && user.props) {
        username = user.props.gitlab_user;
    }

    return {
        username,
        gitlabURL: state['plugins-gitlab'].gitlabURL,
    };
}

export default connect(mapStateToProps)(UserAttribute);
