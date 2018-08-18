import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {
    getAssignedIssues,
    getAssignedMergeRequests,
    getCreatedMergeRequests,
    getUnreads,
} from '../../actions';

import SidebarButtons from './sidebar_buttons.jsx';

function mapStateToProps(state, ownProps) {
    return {
        connected: state['plugins-gitlab'].connected,
        gitlabURL: state['plugins-gitlab'].gitlabURL,
        gitlabUsername: state['plugins-gitlab'].username,
        assignedIssues: state['plugins-gitlab'].assignedIssues,
        assignedMergeRequests: state['plugins-gitlab'].assignedMergeRequests,
        createdMergeRequests: state['plugins-gitlab'].createdMergeRequests,
        // unreads: state['plugins-gitlab'].unreads,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getAssignedIssues,
            getAssignedMergeRequests,
            getCreatedMergeRequests,
            // getUnreads,
        }, dispatch)
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarButtons);
