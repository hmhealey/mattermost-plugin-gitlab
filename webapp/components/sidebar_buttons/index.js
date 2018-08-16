import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getReviews, getUnreads, getYourPrs, getYourAssignments} from '../../actions';

import SidebarButtons from './sidebar_buttons.jsx';

function mapStateToProps(state, ownProps) {
    return {
        connected: state['plugins-gitlab'].connected,
        username: state['plugins-gitlab'].username,
        clientId: state['plugins-gitlab'].clientId,
        // reviews: state['plugins-gitlab'].reviews,
        // yourPrs: state['plugins-gitlab'].yourPrs,
        // yourAssignments: state['plugins-gitlab'].yourAssignments,
        // unreads: state['plugins-gitlab'].unreads,
        // enterpriseURL: state['plugins-gitlab'].enterpriseURL,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            // getReviews,
            // getUnreads,
            // getYourPrs,
            // getYourAssignments,
        }, dispatch)
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarButtons);
