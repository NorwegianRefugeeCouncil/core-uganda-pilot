import PaginationPage from '../pages/pagination.page';
import '../support/commands';

describe('Pagination Page', function () {
    before('Login', () => {
        cy.login('courtney.lare@email.com');
    });
    describe('Should Navigate to Next and Previous Page', () => {
        it('should start on first page', () => {
            const paginationPage = new PaginationPage();
            paginationPage.visitPage('?page=1&perPage=2').getPaginationState().should('contain.text', '1 of');
        });
        it('should navigate to next page', () => {
            const paginationPage = new PaginationPage();
            paginationPage
                .visitPage('?page=1&perPage=2')
                .selectNextPage()
                .getPaginationState()
                .should('contain.text', '2 of');
        });
        it('should navigate to previous page', () => {
            const paginationPage = new PaginationPage();
            paginationPage
                .visitPage('?page=1&perPage=2')
                .selectNextPage()
                .selectPrevPage()
                .getPaginationState()
                .should('contain.text', '1 of');
        });
    });

    describe('Should Navigate to Last and First Page', () => {
        it('should navigate to Last page', () => {
            const paginationPage = new PaginationPage();
            paginationPage
                .visitPage('?page=1&perPage=2')
                .selectLastPage()
                .getPaginationState()
                .should('contain.text', ' of ');
        });
        it('should navigate to First page', () => {
            const paginationPage = new PaginationPage();
            paginationPage
                .visitPage('?page=1&perPage=2')
                .selectNextPage()
                .selectFirstPage()
                .getPaginationState()
                .should('contain.text', '1 of');
        });
    });
});
