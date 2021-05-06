// function stubs
function tintColor() {}
function shadeColor() {}


const vars = {}

class BsVanilla {
  // Variables
//
// Variables should follow the `componentStatePropertySize` formula for
// consistent naming. Ex: this.navLinkDisabledColor and modalContentBoxShadowXs.

// Color system

// scssDocsStart grayColorVariables
public white = "#fff";
public gray100 = "#f8f9fa";
public gray200 = "#e9ecef";
public gray300 = "#dee2e6";
public gray400 = "#ced4da";
public gray500 = "#adb5bd";
public gray600 = "#6c757d";
public gray700 = "#495057";
public gray800 = "#343a40";
public gray900 = "#212529";
public black = "#000";
// scssDocsEnd grayColorVariables

// fusvDisable
// scssDocsStart grayColorsMap
public grays = {
  "100": this.gray100,
  "200": this.gray200,
  "300": this.gray300,
  "400": this.gray400,
  "500": this.gray500,
  "600": this.gray600,
  "700": this.gray700,
  "800": this.gray800,
  "900": this.gray900
};
// scssDocsEnd grayColorsMap
// fusvEnable

// scssDocsStart colorVariables
public blue = "#0d6efd";
public indigo = "#6610f2";
public purple = "#6f42c1";
public pink = "#d63384";
public red = "#dc3545";
public orange = "#fd7e14";
public yellow = "#ffc107";
public green = "#198754";
public teal = "#20c997";
public cyan = "#0dcaf0";
// scssDocsEnd colorVariables

// scssDocsStart colorsMap
public colors = {
  "blue":       this.blue,
  "indigo":     this.indigo,
  "purple":     this.purple,
  "pink":       this.pink,
  "red":        this.red,
  "orange":     this.orange,
  "yellow":     this.yellow,
  "green":      this.green,
  "teal":       this.teal,
  "cyan":       this.cyan,
  "white":      this.white,
  "gray":       this.gray600,
  "grayDark":   this.gray800
};
// scssDocsEnd colorsMap

// scssDocsStart themeColorVariables
public primary = this.blue;
public secondary = this.gray600;
public success = this.green;
public info = this.cyan;
public warning = this.yellow;
public danger = this.red;
public light = this.gray100;
public dark = this.gray900;
// scssDocsEnd themeColorVariables

// scssDocsStart themeColorsMap
public themeColors = {
  "primary":    this.primary,
  "secondary":  this.secondary,
  "success":    this.success,
  "info":       this.info,
  "warning":    this.warning,
  "danger":     this.danger,
  "light":      this.light,
  "dark":       this.dark
};
// scssDocsEnd themeColorsMap

// The contrast ratio to reach against white, to determine if color changes from "light" to "dark". Acceptable values for WCAG 2.0 are 3, 4.5 and 7.
// See https://www.w3.org/TR/WCAG20/#visualAudioContrastContrast
public minContrastRatio = 4.5;

// Customize the light and dark text colors for use in our color contrast function.
public colorContrastDark = this.black;
public colorContrastLight = this.white;

// fusvDisable
public blue100 = tintColor(this.blue, "80%");
public blue200 = tintColor(this.blue, "60%");
public blue300 = tintColor(this.blue, "40%");
public blue400 = tintColor(this.blue, "20%");
public blue500 = this.blue;
public blue600 = shadeColor(this.blue, "20%");
public blue700 = shadeColor(this.blue, "40%");
public blue800 = shadeColor(this.blue, "60%");
public blue900 = shadeColor(this.blue, "80%");

public indigo100 = tintColor(this.indigo, "80%");
public indigo200 = tintColor(this.indigo, "60%");
public indigo300 = tintColor(this.indigo, "40%");
public indigo400 = tintColor(this.indigo, "20%");
public indigo500 = this.indigo;
public indigo600 = shadeColor(this.indigo, "20%");
public indigo700 = shadeColor(this.indigo, "40%");
public indigo800 = shadeColor(this.indigo, "60%");
public indigo900 = shadeColor(this.indigo, "80%");

public purple100 = tintColor(this.purple, "80%");
public purple200 = tintColor(this.purple, "60%");
public purple300 = tintColor(this.purple, "40%");
public purple400 = tintColor(this.purple, "20%");
public purple500 = this.purple;
public purple600 = shadeColor(this.purple, "20%");
public purple700 = shadeColor(this.purple, "40%");
public purple800 = shadeColor(this.purple, "60%");
public purple900 = shadeColor(this.purple, "80%");

public pink100 = tintColor(this.pink, "80%");
public pink200 = tintColor(this.pink, "60%");
public pink300 = tintColor(this.pink, "40%");
public pink400 = tintColor(this.pink, "20%");
public pink500 = this.pink;
public pink600 = shadeColor(this.pink, "20%");
public pink700 = shadeColor(this.pink, "40%");
public pink800 = shadeColor(this.pink, "60%");
public pink900 = shadeColor(this.pink, "80%");

public red100 = tintColor(this.red, "80%");
public red200 = tintColor(this.red, "60%");
public red300 = tintColor(this.red, "40%");
public red400 = tintColor(this.red, "20%");
public red500 = this.red;
public red600 = shadeColor(this.red, "20%");
public red700 = shadeColor(this.red, "40%");
public red800 = shadeColor(this.red, "60%");
public red900 = shadeColor(this.red, "80%");

public orange100 = tintColor(this.orange, "80%");
public orange200 = tintColor(this.orange, "60%");
public orange300 = tintColor(this.orange, "40%");
public orange400 = tintColor(this.orange, "20%");
public orange500 = this.orange;
public orange600 = shadeColor(this.orange, "20%");
public orange700 = shadeColor(this.orange, "40%");
public orange800 = shadeColor(this.orange, "60%");
public orange900 = shadeColor(this.orange, "80%");

public yellow100 = tintColor(this.yellow, "80%");
public yellow200 = tintColor(this.yellow, "60%");
public yellow300 = tintColor(this.yellow, "40%");
public yellow400 = tintColor(this.yellow, "20%");
public yellow500 = this.yellow;
public yellow600 = shadeColor(this.yellow, "20%");
public yellow700 = shadeColor(this.yellow, "40%");
public yellow800 = shadeColor(this.yellow, "60%");
public yellow900 = shadeColor(this.yellow, "80%");

public green100 = tintColor(this.green, "80%");
public green200 = tintColor(this.green, "60%");
public green300 = tintColor(this.green, "40%");
public green400 = tintColor(this.green, "20%");
public green500 = this.green;
public green600 = shadeColor(this.green, "20%");
public green700 = shadeColor(this.green, "40%");
public green800 = shadeColor(this.green, "60%");
public green900 = shadeColor(this.green, "80%");

public teal100 = tintColor(this.teal, "80%");
public teal200 = tintColor(this.teal, "60%");
public teal300 = tintColor(this.teal, "40%");
public teal400 = tintColor(this.teal, "20%");
public teal500 = this.teal;
public teal600 = shadeColor(this.teal, "20%");
public teal700 = shadeColor(this.teal, "40%");
public teal800 = shadeColor(this.teal, "60%");
public teal900 = shadeColor(this.teal, "80%");

public cyan100 = tintColor(this.cyan, "80%");
public cyan200 = tintColor(this.cyan, "60%");
public cyan300 = tintColor(this.cyan, "40%");
public cyan400 = tintColor(this.cyan, "20%");
public cyan500 = this.cyan;
public cyan600 = shadeColor(this.cyan, "20%");
public cyan700 = shadeColor(this.cyan, "40%");
public cyan800 = shadeColor(this.cyan, "60%");
public cyan900 = shadeColor(this.cyan, "80%");
// fusvEnable


// !!!
// Characters which are escaped by the escapeSvg function
// public escapedCharacters = {
//   { "<": "%3c" },
//   { ">": "%3e" },
//   { "#": "%23" },
//   { "(": "%28" },
//   { ")": "%29" },
// };

// Options
//
// Quickly modify global styling by enabling or disabling optional features.

public enableCaret = true;
public enableRounded = true;
public enableShadows = false;
public enableGradients = false;
public enableTransitions = true;
public enableReducedMotion = true;
public enableSmoothScroll = true;
public enableGridClasses = true;
public enableButtonPointers = true;
public enableRfs = true;
public enableValidationIcons = true;
public enableNegativeMargins = false;
public enableDeprecationMessages = true;
public enableImportantUtilities = true;

// Prefix for :root CSS variables

public variablePrefix = "bs-";

// Gradient
//
// The gradient which is added to components if `enableGradients` is `true`
// This gradient is also added to elements with `.bgGradient`
// scssDocsStart variableGradient
public gradient = linearGradient("180deg", rgba(this.white, .15), rgba(this.white, 0));
// scssDocsEnd variableGradient

// Spacing
//
// Control the default styling of most Bootstrap elements by modifying these
// variables. Mostly focused on spacing.
// You can add more entries to the spacers map, should you need more variation.

// scssDocsStart spacerVariablesMaps
public spacer = "1rem";
public spacers = {
  0:  0,
  1:  this.spacer / 4,
  2:  this.spacer / 2,
  3:  this.spacer,
  4:  this.spacer * 1.5,
  5:  this.spacer * 3,
};

// !!!
public negativeSpacers = this.enableNegativeMargins ? negativifyMap(this.spacers) : null;
// scssDocsEnd spacerVariablesMaps

// Position
//
// Define the edge positioning anchors of the position utilities.

// scssDocsStart positionMap
public positionValues = {
  0:  0,
  50:  "50%",
  100: "100%"
};
// scssDocsEnd positionMap

// Body
//
// Settings for the `<body>` element.

public bodyBg = this.white;
public bodyColor = this.gray900;
public bodyTextAlign = null;


// Links
//
// Style anchor elements.

public linkColor = this.primary;
public linkDecoration = "underline";
public linkShadePercentage = "20%";
public linkHoverColor = shiftColor(this.linkColor, this.linkShadePercentage);
public linkHoverDecoration = null;

public stretchedLinkPseudoElement = "after";
public stretchedLinkZIndex:                  1;

// Paragraphs
//
// Style p element.

public paragraphMarginBottom = "1rem";


// Grid breakpoints
//
// Define the minimum dimensions at which your layout will change,
// adapting to different screen sizes, for use in media queries.

// scssDocsStart gridBreakpoints
public gridBreakpoints = {
  xs:  0,
  sm:  "576px",
  md:  "768px",
  lg:  "992px",
  xl:  "1200px",
  xxl = "1400px"
};
// scssDocsEnd gridBreakpoints


// !!!
// @include _assertAscending(this.gridBreakpoints, "gridBreakpoints");
// @include _assertStartsAtZero(this.gridBreakpoints, "gridBreakpoints");


// Grid containers
//
// Define the maximum width of `.container` for different screen sizes.

// scssDocsStart containerMaxWidths
public containerMaxWidths = {
  sm:  "540px",
  md:  "720px",
  lg:  "960px",
  xl:  "1140px",
  xxl = "1320px"
};
// scssDocsEnd containerMaxWidths

// !!! @include _assertAscending(this.containerMaxWidths, "containerMaxWidths");


// Grid columns
//
// Set the number of columns and specify the width of the gutters.

public gridColumns = 12;
public gridGutterWidth = "1.5rem";
public gridRowColumns = 6;

public gutters = this.spacers;

// Container padding

public containerPaddingX = this.gridGutterWidth / 2;


// Components
//
// Define common padding and border radius sizes and more.

// scssDocsStart borderVariables
public borderWidth = "1px";
public borderWidths = {
  1:  "1px",
  2:  "2px",
  3:  "3px",
  4:  "4px",
  5: "5px"
};

public borderColor = this.gray300;
// scssDocsEnd borderVariables

// scssDocsStart borderRadiusVariables
public borderRadius = ".25rem";
public borderRadiusSm = ".2rem";
public borderRadiusLg = ".3rem";
public borderRadiusPill = "50rem";
// scssDocsEnd borderRadiusVariables

// scssDocsStart boxShadowVariables
public boxShadow = `0 .5rem 1rem ${rgba(this.black, .15)}`;
public boxShadowSm = `0 .125rem .25rem ${rgba(this.black, .075)}`;
public boxShadowLg = `0 1rem 3rem ${rgba(this.black, .175)}`;
public boxShadowInset = `${this.inset} 0 1px 2px ${rgba(this.black, .075)}`;
// scssDocsEnd boxShadowVariables

public componentActiveColor = this.white;
public componentActiveBg = this.primary;

// scssDocsStart caretVariables
public caretWidth = ".3em"
public caretVerticalAlign = this.caretWidth * .85;
public caretSpacing = this.caretWidth * .85;
// scssDocsEnd caretVariables

public transitionBase = `${this.all} .2s ease-in-out`;
public transitionFade = `${this.opacity} .15s linear`;
// scssDocsStart collapseTransition
public transitionCollapse = `${this.height} .35s ease`;
// scssDocsEnd collapseTransition

// stylelintDisable functionDisallowedList
// scssDocsStart aspectRatios
public aspectRatios = {
  "1x1": "100%",
  "4x3": this.calc(3 / 4 * "100%"),
  "16x9": this.calc(9 / 16 * "100%"),
  "21x9": this.calc(9 / 21 * "100%")
};
// scssDocsEnd aspectRatios
// stylelintEnable functionDisallowedList

// Typography
//
// Font, lineHeight, and color for body text, headings, and more.

// scssDocsStart fontVariables
// stylelintDisable valueKeywordCase
public fontFamilySansSerif = `${this.systemUi}, -appleSystem, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", "Liberation Sans", sansSerif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"`;
public fontFamilyMonospace =`${this.SFMonoRegular}, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace`;
// stylelintEnable valueKeywordCase
// !!!
public fontFamilyBase = this.var(--#{variablePrefix}fontSansSerif);
public fontFamilyCode = this.var(--#{variablePrefix}fontMonospace);

// fontSizeRoot effects the value of `rem`, which is used for as well font sizes, paddings and margins
// fontSizeBase effects the font size of the body text
public fontSizeRoot = null;
public fontSizeBase = "1rem"; // Assumes the browser default, typically `16px`
public fontSizeSm = this.fontSizeBase * .875;
public fontSizeLg = this.fontSizeBase * 1.25;

public fontWeightLighter = this.lighter;
public fontWeightLight = 300;
public fontWeightNormal = 400;
public fontWeightBold = 700;
public fontWeightBolder = this.bolder;

public fontWeightBase = this.fontWeightNormal;

public lineHeightBase = 1.5;
public lineHeightSm = 1.25;
public lineHeightLg = 2;

public h1FontSize = this.fontSizeBase * 2.5;
public h2FontSize = this.fontSizeBase * 2;
public h3FontSize = this.fontSizeBase * 1.75;
public h4FontSize = this.fontSizeBase * 1.5;
public h5FontSize = this.fontSizeBase * 1.25;
public h6FontSize = this.fontSizeBase;
// scssDocsEnd fontVariables

// scssDocsStart fontSizes
public fontSizes = {
  1:  this.h1FontSize,
  2:  this.h2FontSize,
  3:  this.h3FontSize,
  4:  this.h4FontSize,
  5:  this.h5FontSize,
  6: this.h6FontSize
};
// scssDocsEnd fontSizes

// scssDocsStart headingsVariables
public headingsMarginBottom = this.spacer / 2;
public headingsFontFamily = null;
public headingsFontStyle = null;
public headingsFontWeight = 500;
public headingsLineHeight = 1.2;
public headingsColor = null;
// scssDocsEnd headingsVariables

// scssDocsStart displayHeadings
public displayFontSizes = {
  1:  "5rem",
  2:  "4.5rem",
  3:  "4rem",
  4:  "3.5rem",
  5:  "3rem",
  6: "2.5rem"
};

public displayFontWeight = 300;
public displayLineHeight = this.headingsLineHeight;
// scssDocsEnd displayHeadings

// scssDocsStart typeVariables
public leadFontSize = this.fontSizeBase * 1.25;
public leadFontWeight = 300;

public smallFontSize = ".875em"

public subSupFontSize = ".75em"

public textMuted = this.gray600;

public initialismFontSize = this.smallFontSize;

public blockquoteMarginY = this.spacer;
public blockquoteFontSize = this.fontSizeBase * 1.25;
public blockquoteFooterColor = this.gray600;
public blockquoteFooterFontSize = this.smallFontSize;

public hrMarginY = this.spacer;
public hrColor = this.inherit;
public hrHeight = this.borderWidth;
public hrOpacity = .25;

public legendMarginBottom = ".5rem";
public legendFontSize = "1.5rem";
public legendFontWeight = null;

public markPadding = ".2em"

public dtFontWeight = this.fontWeightBold;

public nestedKbdFontWeight = this.fontWeightBold;

public listInlinePadding = ".5rem";

public markBg = "#fcf8e3";
// scssDocsEnd typeVariables


// Tables
//
// Customizes the `.table` component with basic values, each used across all table variations.

// scssDocsStart tableVariables
public tableCellPaddingY = ".5rem";
public tableCellPaddingX = ".5rem";
public tableCellPaddingYSm =     ".25rem";
public tableCellPaddingXSm =     ".25rem";

public tableCellVerticalAlign = this.top;

public tableColor = this.bodyColor;
public tableBg = this.transparent;

public tableThFontWeight = null;

public tableStripedColor = this.tableColor;
public tableStripedBgFactor = .05;
public tableStripedBg = rgba(this.black, tableStripedBgFactor);

public tableActiveColor = this.tableColor;
public tableActiveBgFactor = .1;
public tableActiveBg = rgba(this.black, tableActiveBgFactor);

public tableHoverColor = this.tableColor;
public tableHoverBgFactor = .075;
public tableHoverBg = rgba(this.black, tableHoverBgFactor);

public tableBorderFactor = .1;
public tableBorderWidth = this.borderWidth;
public tableBorderColor = this.borderColor;

public tableStripedOrder = this.odd;

public tableGroupSeparatorColor = this.currentColor;

public tableCaptionColor = this.textMuted;

public tableBgScale = "-80%";
// scssDocsEnd tableVariables

// scssDocsStart tableLoop
public tableVariants = {
  "primary":    shiftColor(this.primary, tableBgScale),
  "secondary":  shiftColor(this.secondary, tableBgScale),
  "success":    shiftColor(this.success, tableBgScale),
  "info":       shiftColor(this.info, tableBgScale),
  "warning":    shiftColor(this.warning, tableBgScale),
  "danger":     shiftColor(this.danger, tableBgScale),
  "light":      light,
  "dark":       dark,
};
// scssDocsEnd tableLoop


// Buttons + Forms
//
// Shared variables that are reassigned to `input-` and `btn-` specific variables.

// scssDocsStart inputBtnVariables
public inputBtnPaddingY = ".375rem";
public inputBtnPaddingX = ".75rem";
public inputBtnFontFamily = null;
public inputBtnFontSize = this.fontSizeBase;
public inputBtnLineHeight = this.lineHeightBase;

public inputBtnFocusWidth = ".25rem";
public inputBtnFocusColorOpacity = .25;
public inputBtnFocusColor = rgba(this.componentActiveBg, inputBtnFocusColorOpacity);
public inputBtnFocusBlur = 0;
public inputBtnFocusBoxShadow = `0 0 ${this.inputBtnFocusBlur} ${this.inputBtnFocusWidth} ${this.inputBtnFocusColor}`;

public inputBtnPaddingYSm:      ".25rem";
public inputBtnPaddingXSm:      ".5rem";
public inputBtnFontSizeSm = this.fontSizeSm;

public inputBtnPaddingYLg:      ".5rem";
public inputBtnPaddingXLg:      "1rem";
public inputBtnFontSizeLg = this.fontSizeLg;

public inputBtnBorderWidth = this.borderWidth;
// scssDocsEnd inputBtnVariables


// Buttons
//
// For each of Bootstrap's buttons, define text, background, and border color.

// scssDocsStart btnVariables
public btnPaddingY = this.inputBtnPaddingY;
public btnPaddingX = this.inputBtnPaddingX;
public btnFontFamily = this.inputBtnFontFamily;
public btnFontSize = this.inputBtnFontSize;
public btnLineHeight = this.inputBtnLineHeight;
public btnWhiteSpace = null; // Set to `nowrap` to prevent text wrapping

public btnPaddingYSm:            inputBtnPaddingYSm;
public btnPaddingXSm:            inputBtnPaddingXSm;
public btnFontSizeSm = this.inputBtnFontSizeSm;

public btnPaddingYLg:            inputBtnPaddingYLg;
public btnPaddingXLg:            inputBtnPaddingXLg;
public btnFontSizeLg = this.inputBtnFontSizeLg;

public btnBorderWidth = this.inputBtnBorderWidth;

public btnFontWeight = this.fontWeightNormal;
public btnBoxShadow = `${this.inset} 0 1px 0 ${rgba(this.white, .15)}, 0 1px 1px ${rgba(this.black, .075)}`;
public btnFocusWidth = this.inputBtnFocusWidth;
public btnFocusBoxShadow = this.inputBtnFocusBoxShadow;
public btnDisabledOpacity = .65;
public btnActiveBoxShadow = `${this.inset} 0 3px 5px ${rgba(this.black, .125)}`;

public btnLinkColor = this.linkColor;
public btnLinkHoverColor = this.linkHoverColor;
public btnLinkDisabledColor = this.gray600;

// Allows for customizing button radius independently from global border radius
public btnBorderRadius = this.borderRadius;
public btnBorderRadiusSm = this.borderRadiusSm;
public btnBorderRadiusLg = this.borderRadiusLg;

public btnTransition = `color .15s ease-in-out, backgroundColor .15s ease-in-out, borderColor .15s ease-in-out, boxShadow .15s ease-in-out`;

public btnHoverBgShadeAmount = "15%";
public btnHoverBgTintAmount = "15%";
public btnHoverBorderShadeAmount = "20%";
public btnHoverBorderTintAmount = "10%";
public btnActiveBgShadeAmount = "20%";
public btnActiveBgTintAmount = "20%";
public btnActiveBorderShadeAmount = "25%";
public btnActiveBorderTintAmount = "10%";
// scssDocsEnd btnVariables


// Forms

// scssDocsStart formTextVariables
public formTextMarginTop = ".25rem";
public formTextFontSize = this.smallFontSize;
public formTextFontStyle = null;
public formTextFontWeight = null;
public formTextColor = this.textMuted;
// scssDocsEnd formTextVariables

// scssDocsStart formLabelVariables
public formLabelMarginBottom = ".5rem";
public formLabelFontSize = null;
public formLabelFontStyle = null;
public formLabelFontWeight = null;
public formLabelColor = null;
// scssDocsEnd formLabelVariables

// scssDocsStart formInputVariables
public inputPaddingY = this.inputBtnPaddingY;
public inputPaddingX = this.inputBtnPaddingX;
public inputFontFamily = this.inputBtnFontFamily;
public inputFontSize = this.inputBtnFontSize;
public inputFontWeight = this.fontWeightBase;
public inputLineHeight = this.inputBtnLineHeight;

public inputPaddingYSm:                    inputBtnPaddingYSm;
public inputPaddingXSm:                    inputBtnPaddingXSm;
public inputFontSizeSm = this.inputBtnFontSizeSm;

public inputPaddingYLg:                    inputBtnPaddingYLg;
public inputPaddingXLg:                    inputBtnPaddingXLg;
public inputFontSizeLg = this.inputBtnFontSizeLg;

public inputBg = this.white;
public inputDisabledBg = this.gray200;
public inputDisabledBorderColor = null;

public inputColor = this.bodyColor;
public inputBorderColor = this.gray400;
public inputBorderWidth = this.inputBtnBorderWidth;
public inputBoxShadow = this.boxShadowInset;

public inputBorderRadius = this.borderRadius;
public inputBorderRadiusSm = this.borderRadiusSm;
public inputBorderRadiusLg = this.borderRadiusLg;

public inputFocusBg = this.inputBg;
public inputFocusBorderColor = tintColor(this.componentActiveBg, "50%");
public inputFocusColor = this.inputColor;
public inputFocusWidth = this.inputBtnFocusWidth;
public inputFocusBoxShadow = this.inputBtnFocusBoxShadow;

public inputPlaceholderColor = this.gray600;
public inputPlaintextColor = this.bodyColor;

public inputHeightBorder = this.inputBorderWidth * 2;

public inputHeightInner = this.add(this.inputLineHeight * "1em", inputPaddingY * 2);
public inputHeightInnerHalf = this.add(this.inputLineHeight * ".5em", inputPaddingY);
public inputHeightInnerQuarter = this.add(this.inputLineHeight * ".25em", inputPaddingY / 2);

public inputHeight = this.add(this.inputLineHeight * "1em", add(this.inputPaddingY * 2, inputHeightBorder, false));
public inputHeightSm = this.add(this.inputLineHeight * "1em", add(this.inputPaddingYSm * 2, inputHeightBorder, false));
public inputHeightLg = this.add(this.inputLineHeight * "1em", add(this.inputPaddingYLg * 2, inputHeightBorder, false));

public inputTransition = `${this.borderColor} .15s ease-in-out, boxShadow .15s ease-in-out`;
// scssDocsEnd formInputVariables

// scssDocsStart formCheckVariables
public formCheckInputWidth = "1em"
public formCheckMinHeight = this.fontSizeBase * lineHeightBase;
public formCheckPaddingStart = this.formCheckInputWidth + ".5em"
public formCheckMarginBottom = ".125rem";
public formCheckLabelColor = null;
public formCheckLabelCursor = null;
public formCheckTransition = null;

public formCheckInputActiveFilter = this.brightness("90%");

public formCheckInputBg = this.inputBg;
public formCheckInputBorder = `1px solid ${rgba(0, 0, 0, .25)}`;
public formCheckInputBorderRadius = ".25em"
public formCheckRadioBorderRadius = "50%";
public formCheckInputFocusBorder = this.inputFocusBorderColor;
public formCheckInputFocusBoxShadow = this.inputBtnFocusBoxShadow;

public formCheckInputCheckedColor = this.componentActiveColor;
public formCheckInputCheckedBgColor = this.componentActiveBg;
public formCheckInputCheckedBorderColor = this.formCheckInputCheckedBgColor;
public formCheckInputCheckedBgImage = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20'><path fill='none' stroke='#{formCheckInputCheckedColor}' strokeLinecap='round' strokeLinejoin='round' strokeWidth='3' d='M6 10l3 3l66'/></svg>");
public formCheckRadioCheckedBgImage = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='-4 -4 8 8'><circle r='2' fill='#{formCheckInputCheckedColor}'/></svg>");

public formCheckInputIndeterminateColor = this.componentActiveColor;
public formCheckInputIndeterminateBgColor = this.componentActiveBg;
public formCheckInputIndeterminateBorderColor = this.formCheckInputIndeterminateBgColor;
public formCheckInputIndeterminateBgImage = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20'><path fill='none' stroke='#{formCheckInputIndeterminateColor}' strokeLinecap='round' strokeLinejoin='round' strokeWidth='3' d='M6 10h8'/></svg>");

public formCheckInputDisabledOpacity = .5;
public formCheckLabelDisabledOpacity = this.formCheckInputDisabledOpacity;
public formCheckBtnCheckDisabledOpacity = this.btnDisabledOpacity;

public formCheckInlineMarginEnd = "1rem";
// scssDocsEnd formCheckVariables

// scssDocsStart formSwitchVariables
public formSwitchColor = rgba(0, 0, 0, .25);
public formSwitchWidth = "2em"
public formSwitchPaddingStart = this.formSwitchWidth + ".5em"
public formSwitchBgImage = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='-4 -4 8 8'><circle r='3' fill='#{formSwitchColor}'/></svg>");
public formSwitchBorderRadius = this.formSwitchWidth;
public formSwitchTransition = `${this.backgroundPosition} .15s ease-in-out`;

public formSwitchFocusColor = this.inputFocusBorderColor;
public formSwitchFocusBgImage = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='-4 -4 8 8'><circle r='3' fill='#{formSwitchFocusColor}'/></svg>");

public formSwitchCheckedColor = this.componentActiveColor;
public formSwitchCheckedBgImage = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='-4 -4 8 8'><circle r='3' fill='#{formSwitchCheckedColor}'/></svg>");
public formSwitchCheckedBgPosition = "right center";
// scssDocsEnd formSwitchVariables

// scssDocsStart inputGroupVariables
public inputGroupAddonPaddingY = this.inputPaddingY;
public inputGroupAddonPaddingX = this.inputPaddingX;
public inputGroupAddonFontWeight = this.inputFontWeight;
public inputGroupAddonColor = this.inputColor;
public inputGroupAddonBg = this.gray200;
public inputGroupAddonBorderColor = this.inputBorderColor;
// scssDocsEnd inputGroupVariables

// scssDocsStart formSelectVariables
public formSelectPaddingY = this.inputPaddingY;
public formSelectPaddingX = this.inputPaddingX;
public formSelectFontFamily = this.inputFontFamily;
public formSelectFontSize = this.inputFontSize;
public formSelectIndicatorPadding = this.formSelectPaddingX * 3; // Extra padding for backgroundImage
public formSelectFontWeight = this.inputFontWeight;
public formSelectLineHeight = this.inputLineHeight;
public formSelectColor = this.inputColor;
public formSelectBg = this.inputBg;
public formSelectDisabledColor = null;
public formSelectDisabledBg = this.gray200;
public formSelectDisabledBorderColor = this.inputDisabledBorderColor;
public formSelectBgPosition = `this.right formSelectPaddingX center`;
public formSelectBgSize = "16px 12px"; // In pixels because image dimensions
public formSelectIndicatorColor = this.gray800;
public formSelectIndicator = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'><path fill='none' stroke='#{formSelectIndicatorColor}' strokeLinecap='round' strokeLinejoin='round' strokeWidth='2' d='M2 5l6 6 66'/></svg>");

public formSelectFeedbackIconPaddingEnd = `${this.formSelectPaddingX} * 2.5 + ${this.formSelectIndicatorPadding}`;
public formSelectFeedbackIconPosition = "center right ${this.formSelectIndicatorPadding}";
public formSelectFeedbackIconSize = `${this.inputHeightInnerHalf} ${this.inputHeightInnerHalf}`;

public formSelectBorderWidth = this.inputBorderWidth;
public formSelectBorderColor = this.inputBorderColor;
public formSelectBorderRadius = this.borderRadius;
public formSelectBoxShadow = this.boxShadowInset;

public formSelectFocusBorderColor = this.inputFocusBorderColor;
public formSelectFocusWidth = this.inputFocusWidth;
public formSelectFocusBoxShadow = `0 0 0 ${this.formSelectFocusWidth} ${this.inputBtnFocusColor}`;

public formSelectPaddingYSm:        inputPaddingYSm;
public formSelectPaddingXSm:        inputPaddingXSm;
public formSelectFontSizeSm = this.inputFontSizeSm;

public formSelectPaddingYLg:        inputPaddingYLg;
public formSelectPaddingXLg:        inputPaddingXLg;
public formSelectFontSizeLg = this.inputFontSizeLg;
// scssDocsEnd formSelectVariables

// scssDocsStart formRangeVariables
public formRangeTrackWidth = "100%";
public formRangeTrackHeight = ".5rem";
public formRangeTrackCursor = this.pointer;
public formRangeTrackBg = this.gray300;
public formRangeTrackBorderRadius = "1rem";
public formRangeTrackBoxShadow = this.boxShadowInset;

public formRangeThumbWidth = "1rem";
public formRangeThumbHeight = this.formRangeThumbWidth;
public formRangeThumbBg = this.componentActiveBg;
public formRangeThumbBorder = 0;
public formRangeThumbBorderRadius = "1rem";
public formRangeThumbBoxShadow = `0 .1rem .25rem ${rgba(this.black, .1)}`;
public formRangeThumbFocusBoxShadow = `0 0 0 1px ${this.bodyBg}, ${this.inputFocusBoxShadow}`;
public formRangeThumbFocusBoxShadowWidth = this.inputFocusWidth; // For focus box shadow issue in Edge
public formRangeThumbActiveBg = tintColor(this.componentActiveBg, "70%");
public formRangeThumbDisabledBg = this.gray500;
public formRangeThumbTransition = `${this.backgroundColor} .15s ease-in-out, borderColor .15s ease-in-out, boxShadow .15s ease-in-out`;
// scssDocsEnd formRangeVariables

// scssDocsStart formFileVariables
public formFileButtonColor = this.inputColor;
public formFileButtonBg = this.inputGroupAddonBg;
public formFileButtonHoverBg = shadeColor(this.formFileButtonBg, "5%");
// scssDocsEnd formFileVariables

// scssDocsStart formFloatingVariables
public formFloatingHeight = add("3.5rem", this.inputHeightBorder);
public formFloatingPaddingX = this.inputPaddingX;
public formFloatingPaddingY = "1rem";
public formFloatingInputPaddingT = "1.625rem";
public formFloatingInputPaddingB = ".625rem";
public formFloatingLabelOpacity = .65;
public formFloatingLabelTransform = `${scale(.85)} ${translateY("-.5rem")} ${translateX(.15rem)}`;
public formFloatingTransition = `${this.opacity} .1s ease-in-out, transform .1s ease-in-out`;
// scssDocsEnd formFloatingVariables

// Form validation

// scssDocsStart formFeedbackVariables
public formFeedbackMarginTop = this.formTextMarginTop;
public formFeedbackFontSize = this.formTextFontSize;
public formFeedbackFontStyle = this.formTextFontStyle;
public formFeedbackValidColor = this.success;
public formFeedbackInvalidColor = this.danger;

public formFeedbackIconValidColor = this.formFeedbackValidColor;
public formFeedbackIconValid = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8'><path fill='#{formFeedbackIconValidColor}' d='M2.3 6.73L.6 4.53c-.41.04.461.4 1.1-.8l1.1 1.4 3.43.8c.6-.63 1.6-.27 1.2.7l4 4.6c-.43.5-.8.41.1.1z'/></svg>");
public formFeedbackIconInvalidColor = this.formFeedbackInvalidColor;
public formFeedbackIconInvalid = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' width='12' height='12' fill='none' stroke='#{formFeedbackIconInvalidColor}'><circle cx='6' cy='6' r='4.5'/><path strokeLinejoin='round' d='M5.8 3.6h.4L6 6.5z'/><circle cx='6' cy='8.2' r='.6' fill='#{formFeedbackIconInvalidColor}' stroke='none'/></svg>");
// scssDocsEnd formFeedbackVariables

// scssDocsStart formValidationStates
public formValidationStates = {
  "valid": {
    "color": this.formFeedbackValidColor,
    "icon": this.formFeedbackIconValid
  },
  "invalid": {
    "color": this.formFeedbackInvalidColor,
    "icon": this.formFeedbackIconInvalid
  }
};
// scssDocsEnd formValidationStates

// ZIndex master list
//
// Warning: this.Avoid customizing these values. They're used for a bird's eye view
// of components dependent on the zAxis and are designed to all work together.

// scssDocsStart zindexStack
public zindexDropdown = 1000;
public zindexSticky = 1020;
public zindexFixed = 1030;
public zindexOffcanvas = 1040;
public zindexModalBackdrop = 1050;
public zindexModal = 1060;
public zindexPopover = 1070;
public zindexTooltip = 1080;
// scssDocsEnd zindexStack


// Navs

// scssDocsStart navVariables
public navLinkPaddingY = ".5rem";
public navLinkPaddingX = "1rem";
public navLinkFontSize = null;
public navLinkFontWeight = null;
public navLinkColor = null;
public navLinkHoverColor = null;
public navLinkTransition = `color .15s ease-in-out, background-color .15s ease-in-out, border-color .15s ease-in-out`;
public navLinkDisabledColor = this.gray600;

public navTabsBorderColor = this.gray300;
public navTabsBorderWidth = this.borderWidth;
public navTabsBorderRadius = this.borderRadius;
public navTabsLinkHoverBorderColor = `${this.gray200} ${this.gray200} ${this.navTabsBorderColor}`;
public navTabsLinkActiveColor = this.gray700;
public navTabsLinkActiveBg = this.bodyBg;
public navTabsLinkActiveBorderColor = `${this.gray300} ${this.gray300} ${this.navTabsLinkActiveBg}`;

public navPillsBorderRadius = this.borderRadius;
public navPillsLinkActiveColor = this.componentActiveColor;
public navPillsLinkActiveBg = this.componentActiveBg;
// scssDocsEnd navVariables


// Navbar

// scssDocsStart navbarVariables
public navbarPaddingY = this.spacer / 2;
public navbarPaddingX = null;

public navbarNavLinkPaddingX = ".5rem";

public navbarBrandFontSize = this.fontSizeLg;
// Compute the navbarBrand paddingY so the navbarBrand will have the same height as navbarText and navLink
public navLinkHeight = this.fontSizeBase * lineHeightBase + navLinkPaddingY * 2;
public navbarBrandHeight = this.navbarBrandFontSize * lineHeightBase;
public navbarBrandPaddingY = (this.navLinkHeight - navbarBrandHeight) / 2;
public navbarBrandMarginEnd = "1rem";

public navbarTogglerPaddingY = ".25rem";
public navbarTogglerPaddingX = ".75rem";
public navbarTogglerFontSize = this.fontSizeLg;
public navbarTogglerBorderRadius = this.btnBorderRadius;
public navbarTogglerFocusWidth = this.btnFocusWidth;
public navbarTogglerTransition = 'box-shadow .15s ease-in-out';
// scssDocsEnd navbarVariables

// scssDocsStart navbarThemeVariables
public navbarDarkColor = rgba(this.white, .55);
public navbarDarkHoverColor = rgba(this.white, .75);
public navbarDarkActiveColor = this.white;
public navbarDarkDisabledColor = rgba(this.white, .25);
public navbarDarkTogglerIconBg = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 30 30'><path stroke='#{navbarDarkColor}' strokeLinecap='round' strokeMiterlimit='10' strokeWidth='2' d='M4 7h22M4 15h22M4 23h22'/></svg>");
public navbarDarkTogglerBorderColor = rgba(this.white, .1);

public navbarLightColor = rgba(this.black, .55);
public navbarLightHoverColor = rgba(this.black, .7);
public navbarLightActiveColor = rgba(this.black, .9);
public navbarLightDisabledColor = rgba(this.black, .3);
public navbarLightTogglerIconBg = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 30 30'><path stroke='#{navbarLightColor}' strokeLinecap='round' strokeMiterlimit='10' strokeWidth='2' d='M4 7h22M4 15h22M4 23h22'/></svg>");
public navbarLightTogglerBorderColor = rgba(this.black, .1);

public navbarLightBrandColor = this.navbarLightActiveColor;
public navbarLightBrandHoverColor = this.navbarLightActiveColor;
public navbarDarkBrandColor = this.navbarDarkActiveColor;
public navbarDarkBrandHoverColor = this.navbarDarkActiveColor;
// scssDocsEnd navbarThemeVariables


// Dropdowns
//
// Dropdown menu container and contents.

// scssDocsStart dropdownVariables
public dropdownMinWidth = "10rem";
public dropdownPaddingX = 0;
public dropdownPaddingY = ".5rem";
public dropdownSpacer = ".125rem";
public dropdownFontSize = this.fontSizeBase;
public dropdownColor = this.bodyColor;
public dropdownBg = this.white;
public dropdownBorderColor = rgba(this.black, .15);
public dropdownBorderRadius = this.borderRadius;
public dropdownBorderWidth = this.borderWidth;
public dropdownInnerBorderRadius = subtract(this.dropdownBorderRadius, dropdownBorderWidth);
public dropdownDividerBg = this.dropdownBorderColor;
public dropdownDividerMarginY = this.spacer / 2;
public dropdownBoxShadow = this.boxShadow;

public dropdownLinkColor = this.gray900;
public dropdownLinkHoverColor = shadeColor(this.gray900, "10%");
public dropdownLinkHoverBg = this.gray200;

public dropdownLinkActiveColor = this.componentActiveColor;
public dropdownLinkActiveBg = this.componentActiveBg;

public dropdownLinkDisabledColor = this.gray500;

public dropdownItemPaddingY = this.spacer / 4;
public dropdownItemPaddingX = this.spacer;

public dropdownHeaderColor = this.gray600;
public dropdownHeaderPadding = `${this.dropdownPaddingY} ${this.dropdownItemPaddingX}`;
// scssDocsEnd dropdownVariables

// scssDocsStart dropdownDarkVariables
public dropdownDarkColor = this.gray300;
public dropdownDarkBg = this.gray800;
public dropdownDarkBorderColor = this.dropdownBorderColor;
public dropdownDarkDividerBg = this.dropdownDividerBg;
public dropdownDarkBoxShadow = null;
public dropdownDarkLinkColor = this.dropdownDarkColor;
public dropdownDarkLinkHoverColor = this.white;
public dropdownDarkLinkHoverBg = rgba(this.white, .15);
public dropdownDarkLinkActiveColor = this.dropdownLinkActiveColor;
public dropdownDarkLinkActiveBg = this.dropdownLinkActiveBg;
public dropdownDarkLinkDisabledColor = this.gray500;
public dropdownDarkHeaderColor = this.gray500;
// scssDocsEnd dropdownDarkVariables


// Pagination

// scssDocsStart paginationVariables
public paginationPaddingY = ".375rem";
public paginationPaddingX = ".75rem";
public paginationPaddingYSm = ".25rem";
public paginationPaddingXSm = ".5rem";
public paginationPaddingYLg = ".75rem";
public paginationPaddingXLg = "1.5rem";

public paginationColor = this.linkColor;
public paginationBg = this.white;
public paginationBorderWidth = this.borderWidth;
public paginationBorderRadius = this.borderRadius;
public paginationMarginStart = -paginationBorderWidth;
public paginationBorderColor = this.gray300;

public paginationFocusColor = this.linkHoverColor;
public paginationFocusBg = this.gray200;
public paginationFocusBoxShadow = this.inputBtnFocusBoxShadow;
public paginationFocusOutline = 0;

public paginationHoverColor = this.linkHoverColor;
public paginationHoverBg = this.gray200;
public paginationHoverBorderColor = this.gray300;

public paginationActiveColor = this.componentActiveColor;
public paginationActiveBg = this.componentActiveBg;
public paginationActiveBorderColor = this.paginationActiveBg;

public paginationDisabledColor = this.gray600;
public paginationDisabledBg = this.white;
public paginationDisabledBorderColor = this.gray300;

public paginationTransition = 'color .15s ease-in-out, background-olor .15s ease-in-out, border-color .15s ease-in-out, box-shadow .15s ease-in-out';

public paginationBorderRadiusSm = this.borderRadiusSm;
public paginationBorderRadiusLg = this.borderRadiusLg;
// scssDocsEnd paginationVariables


// Cards

// scssDocsStart cardVariables
public cardSpacerY = this.spacer;
public cardSpacerX = this.spacer;
public cardTitleSpacerY = this.spacer / 2;
public cardBorderWidth = this.borderWidth;
public cardBorderRadius = this.borderRadius;
public cardBorderColor = rgba(this.black, .125);
public cardInnerBorderRadius = subtract(this.cardBorderRadius, cardBorderWidth);
public cardCapPaddingY = this.cardSpacerY / 2;
public cardCapPaddingX = this.cardSpacerX;
public cardCapBg = rgba(this.black, .03);
public cardCapColor = null;
public cardHeight = null;
public cardColor = null;
public cardBg = this.white;
public cardImgOverlayPadding = this.spacer;
public cardGroupMargin = this.gridGutterWidth / 2;
// scssDocsEnd cardVariables

// Accordion

// scssDocsStart accordionVariables
public accordionPaddingY = "1rem";
public accordionPaddingX = "1.25rem";
public accordionColor = this.bodyColor;
public accordionBg = this.bodyBg;
public accordionBorderWidth = this.borderWidth;
public accordionBorderColor = rgba(this.black, .125);
public accordionBorderRadius = this.borderRadius;
public accordionInnerBorderRadius = subtract(this.accordionBorderRadius, accordionBorderWidth);

public accordionBodyPaddingY = this.accordionPaddingY;
public accordionBodyPaddingX = this.accordionPaddingX;

public accordionButtonPaddingY = this.accordionPaddingY;
public accordionButtonPaddingX = this.accordionPaddingX;
public accordionButtonColor = this.accordionColor;
public accordionButtonBg = this.accordionBg;
public accordionTransition = `${this.btnTransition}, borde-radius .15s ease`;
public accordionButtonActiveBg = tintColor(this.componentActiveBg, "90%");
public accordionButtonActiveColor = shadeColor(this.primary, "10%");

public accordionButtonFocusBorderColor = this.inputFocusBorderColor;
public accordionButtonFocusBoxShadow = this.btnFocusBoxShadow;

public accordionIconWidth = "1.25rem";
public accordionIconColor = this.accordionColor;
public accordionIconActiveColor = this.accordionButtonActiveColor;
public accordionIconTransition = 'transform .2s ease-in-out';
public accordionIconTransform = rotate("180deg");

public accordionButtonIcon = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='#{accordionIconColor}'><path fillRule='evenodd' d='M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.6465.647a.5.5 0 0 1 .708.708l6 6a.5.5 0 0 1-.708 0l66a.5.5 0 0 1 0-.708z'/></svg>");
public accordionButtonActiveIcon = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='#{accordionIconActiveColor}'><path fillRule='evenodd' d='M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.6465.647a.5.5 0 0 1 .708.708l6 6a.5.5 0 0 1-.708 0l66a.5.5 0 0 1 0-.708z'/></svg>");
// scssDocsEnd accordionVariables

// Tooltips

// scssDocsStart tooltipVariables
public tooltipFontSize = this.fontSizeSm;
public tooltipMaxWidth = "200px";
public tooltipColor = this.white;
public tooltipBg = this.black;
public tooltipBorderRadius = this.borderRadius;
public tooltipOpacity = .9;
public tooltipPaddingY = this.spacer / 4;
public tooltipPaddingX = this.spacer / 2;
public tooltipMargin = 0;

public tooltipArrowWidth = ".8rem";
public tooltipArrowHeight = ".4rem";
public tooltipArrowColor = this.tooltipBg;
// scssDocsEnd tooltipVariables

// Form tooltips must come after regular tooltips
// scssDocsStart tooltipFeedbackVariables
public formFeedbackTooltipPaddingY = this.tooltipPaddingY;
public formFeedbackTooltipPaddingX = this.tooltipPaddingX;
public formFeedbackTooltipFontSize = this.tooltipFontSize;
public formFeedbackTooltipLineHeight = null;
public formFeedbackTooltipOpacity = this.tooltipOpacity;
public formFeedbackTooltipBorderRadius = this.tooltipBorderRadius;
// scssDocsStart tooltipFeedbackVariables


// Popovers

// scssDocsStart popoverVariables
public popoverFontSize = this.fontSizeSm;
public popoverBg = this.white;
public popoverMaxWidth = "276px";
public popoverBorderWidth = this.borderWidth;
public popoverBorderColor = rgba(this.black, .2);
public popoverBorderRadius = this.borderRadiusLg;
public popoverInnerBorderRadius = subtract(this.popoverBorderRadius, popoverBorderWidth);
public popoverBoxShadow = this.boxShadow;

public popoverHeaderBg = shadeColor(this.popoverBg, "6%");
public popoverHeaderColor = this.headingsColor;
public popoverHeaderPaddingY = ".5rem";
public popoverHeaderPaddingX = this.spacer;

public popoverBodyColor = this.bodyColor;
public popoverBodyPaddingY = this.spacer;
public popoverBodyPaddingX = this.spacer;

public popoverArrowWidth = "1rem";
public popoverArrowHeight = ".5rem";
public popoverArrowColor = this.popoverBg;

public popoverArrowOuterColor = this.fadeIn(this.popoverBorderColor, .05);
// scssDocsEnd popoverVariables


// Toasts

// scssDocsStart toastVariables
public toastMaxWidth = "350px";
public toastPaddingX = ".75rem";
public toastPaddingY = ".5rem";
public toastFontSize = ".875rem";
public toastColor = null;
public toastBackgroundColor = rgba(this.white, .85);
public toastBorderWidth = "1px";
public toastBorderColor = rgba(0, 0, 0, .1);
public toastBorderRadius = this.borderRadius;
public toastBoxShadow = this.boxShadow;
public toastSpacing = this.containerPaddingX;

public toastHeaderColor = this.gray600;
public toastHeaderBackgroundColor = rgba(this.white, .85);
public toastHeaderBorderColor = rgba(0, 0, 0, .05);
// scssDocsEnd toastVariables


// Badges

// scssDocsStart badgeVariables
public badgeFontSize = ".75em"
public badgeFontWeight = this.fontWeightBold;
public badgeColor = this.white;
public badgePaddingY = ".35em"
public badgePaddingX = ".65em"
public badgeBorderRadius = this.borderRadius;
// scssDocsEnd badgeVariables


// Modals

// scssDocsStart modalVariables
public modalInnerPadding = this.spacer;

public modalFooterMarginBetween = ".5rem";

public modalDialogMargin = ".5rem";
public modalDialogMarginYSmUp:       "1.75rem";

public modalTitleLineHeight = this.lineHeightBase;

public modalContentColor = null;
public modalContentBg = this.white;
public modalContentBorderColor = rgba(this.black, .2);
public modalContentBorderWidth = this.borderWidth;
public modalContentBorderRadius = this.borderRadiusLg;
public modalContentInnerBorderRadius = subtract(this.modalContentBorderRadius, modalContentBorderWidth);
public modalContentBoxShadowXs = this.boxShadowSm;
public modalContentBoxShadowSmUp = this.boxShadow;

public modalBackdropBg = this.black;
public modalBackdropOpacity = .5;
public modalHeaderBorderColor = this.borderColor;
public modalFooterBorderColor = this.modalHeaderBorderColor;
public modalHeaderBorderWidth = this.modalContentBorderWidth;
public modalFooterBorderWidth = this.modalHeaderBorderWidth;
public modalHeaderPaddingY = this.modalInnerPadding;
public modalHeaderPaddingX = this.modalInnerPadding;
public modalHeaderPadding = `${this.modalHeaderPaddingY} ${this.modalHeaderPaddingX}`; // Keep this for backwards compatibility

public modalSm = "300px";
public modalMd = "500px";
public modalLg = "800px";
public modalXl = "1140px";

public modalFadeTransform = translate(0, "-50px");
public modalShowTransform = "none";
public modalTransition = 'transform .3s ease-out';
public modalScaleTransform = scale(1.02);
// scssDocsEnd modalVariables


// Alerts
//
// Define alert colors, border radius, and padding.

// scssDocsStart alertVariables
public alertPaddingY = this.spacer;
public alertPaddingX = this.spacer;
public alertMarginBottom = "1rem";
public alertBorderRadius = this.borderRadius;
public alertLinkFontWeight = this.fontWeightBold;
public alertBorderWidth = this.borderWidth;
public alertBgScale = "-80%";
public alertBorderScale = "-70%";
public alertColorScale = "40%";
public alertDismissiblePaddingR = this.alertPaddingX * 3; // 3x covers width of x plus default padding on either side
// scssDocsEnd alertVariables


// Progress bars

// scssDocsStart progressVariables
public progressHeight = "1rem";
public progressFontSize = this.fontSizeBase * .75;
public progressBg = this.gray200;
public progressBorderRadius = this.borderRadius;
public progressBoxShadow = this.boxShadowInset;
public progressBarColor = this.white;
public progressBarBg = this.primary;
public progressBarAnimationTiming = '1s linear infinite';
public progressBarTransition = 'width .6s ease';
// scssDocsEnd progressVariables


// List group

// scssDocsStart listGroupVariables
public listGroupColor = this.gray900;
public listGroupBg = this.white;
public listGroupBorderColor = rgba(this.black, .125);
public listGroupBorderWidth = this.borderWidth;
public listGroupBorderRadius = this.borderRadius;

public listGroupItemPaddingY = this.spacer / 2;
public listGroupItemPaddingX = this.spacer;
public listGroupItemBgScale = "-80%";
public listGroupItemColorScale = "40%";

public listGroupHoverBg = this.gray100;
public listGroupActiveColor = this.componentActiveColor;
public listGroupActiveBg = this.componentActiveBg;
public listGroupActiveBorderColor = this.listGroupActiveBg;

public listGroupDisabledColor = this.gray600;
public listGroupDisabledBg = this.listGroupBg;

public listGroupActionColor = this.gray700;
public listGroupActionHoverColor = this.listGroupActionColor;

public listGroupActionActiveColor = this.bodyColor;
public listGroupActionActiveBg = this.gray200;
// scssDocsEnd listGroupVariables


// Image thumbnails

// scssDocsStart thumbnailVariables
public thumbnailPadding = ".25rem";
public thumbnailBg = this.bodyBg;
public thumbnailBorderWidth = this.borderWidth;
public thumbnailBorderColor = this.gray300;
public thumbnailBorderRadius = this.borderRadius;
public thumbnailBoxShadow = this.boxShadowSm;
// scssDocsEnd thumbnailVariables


// Figures

// scssDocsStart figureVariables
public figureCaptionFontSize = this.smallFontSize;
public figureCaptionColor = this.gray600;
// scssDocsEnd figureVariables


// Breadcrumbs

// scssDocsStart breadcrumbVariables
public breadcrumbFontSize = null;
public breadcrumbPaddingY = 0;
public breadcrumbPaddingX = 0;
public breadcrumbItemPaddingX = ".5rem";
public breadcrumbMarginBottom = "1rem";
public breadcrumbBg = null;
public breadcrumbDividerColor = this.gray600;
public breadcrumbActiveColor = this.gray600;
public breadcrumbDivider = this.quote("/");
public breadcrumbDividerFlipped = this.breadcrumbDivider;
public breadcrumbBorderRadius = null;
// scssDocsEnd breadcrumbVariables

// Carousel

// scssDocsStart carouselVariables
public carouselControlColor = this.white;
public carouselControlWidth = "15%";
public carouselControlOpacity = .5;
public carouselControlHoverOpacity = .9;
public carouselControlTransition = 'opacity .15s ease';

public carouselIndicatorWidth = "30px";
public carouselIndicatorHeight = "3px";
public carouselIndicatorHitAreaHeight = "10px";
public carouselIndicatorSpacer = "3px";
public carouselIndicatorOpacity = .5;
public carouselIndicatorActiveBg = this.white;
public carouselIndicatorActiveOpacity = 1;
public carouselIndicatorTransition = 'opacity .6s ease';

public carouselCaptionWidth = "70%";
public carouselCaptionColor = this.white;
public carouselCaptionPaddingY = "1.25rem";
public carouselCaptionSpacer = "1.25rem";

public carouselControlIconWidth = "2rem";

public carouselControlPrevIconBg = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='#{carouselControlColor}'><path d='M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l66a.5.5 0 0 1 0-.708l66a.5.5 0 0 1 .708 0z'/></svg>");
public carouselControlNextIconBg = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='#{carouselControlColor}'><path d='M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z'/></svg>");

public carouselTransitionDuration = ".6s";
public carouselTransition = `transform ${this.carouselTransitionDuration} ease-in-out`; // Define transform transition first if using multiple transitions (this.e.g., `transform 2s ease, opacity .5s ease-out`)

public carouselDarkIndicatorActiveBg = this.black;
public carouselDarkCaptionColor = this.black;
public carouselDarkControlIconFilter = `${invert(1)} ${grayscale(100)}`;
// scssDocsEnd carouselVariables


// Spinners

// scssDocsStart spinnerVariables
public spinnerWidth = "2rem";
public spinnerHeight = this.spinnerWidth;
public spinnerBorderWidth = ".25em"
public spinnerAnimationSpeed = ".75s";

public spinnerWidthSm = "1rem";
public spinnerHeightSm = this.spinnerWidthSm;
public spinnerBorderWidthSm = ".2em"
// scssDocsEnd spinnerVariables


// Close

// scssDocsStart closeVariables
public btnCloseWidth = "1em"
public btnCloseHeight = this.btnCloseWidth;
public btnClosePaddingX = ".25em"
public btnClosePaddingY = this.btnClosePaddingX;
public btnCloseColor = this.black;
public btnCloseBg = url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='#{btnCloseColor}'><path d='M.293.293a1 1 0 011.414 0L8 6.586 14.293.293a1 1 0 111.414 1.414L9.414 8l6.293 6.293a1 1 0 011.414 1.414L8 9.414l6.293 6.293a1 1 0 011.4141.414L6.586 8 .293 1.707a1 1 0 0101.414z'/></svg>");
public btnCloseFocusShadow = this.inputBtnFocusBoxShadow;
public btnCloseOpacity = .5;
public btnCloseHoverOpacity = .75;
public btnCloseFocusOpacity = 1;
public btnCloseDisabledOpacity = .25;
public btnCloseWhiteFilter = `${invert(1)} ${grayscale('100%')} ${brightness('200%')}`;
// scssDocsEnd closeVariables


// Offcanvas

// scssDocsStart offcanvasVariables
public offcanvasPaddingY = this.modalInnerPadding;
public offcanvasPaddingX = this.modalInnerPadding;
public offcanvasHorizontalWidth = "400px";
public offcanvasVerticalHeight = "30vh"
public offcanvasTransitionDuration = ".3s";
public offcanvasBorderColor = this.modalContentBorderColor;
public offcanvasBorderWidth = this.modalContentBorderWidth;
public offcanvasTitleLineHeight = this.modalTitleLineHeight;
public offcanvasBgColor = this.modalContentBg;
public offcanvasColor = this.modalContentColor;
public offcanvasBodyBackdropColor = rgba(this.modalBackdropBg, modalBackdropOpacity);
public offcanvasBoxShadow = this.modalContentBoxShadowXs;
// scssDocsEnd offcanvasVariables

// Code

public codeFontSize = this.smallFontSize;
public codeColor = this.pink;

public kbdPaddingY = ".2rem";
public kbdPaddingX = ".4rem";
public kbdFontSize = this.codeFontSize;
public kbdColor = this.white;
public kbdBg = this.gray900;

public preColor = null;

}


const bsVanillaInstance = new BsVanilla()
export default bsVanillaInstance