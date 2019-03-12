import pandas as pd


def prepare_dataframe(data) -> pd.DataFrame:
    df = pd.DataFrame.from_dict(data)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
    df.set_index('timestamp', inplace=True)
    df.sort_index(inplace=True)
    return df


def post_generate(plt):
    plt.spines['right'].set_visible(False)
    plt.spines['top'].set_visible(False)
    plt.spines['left'].set_visible(False)
    plt.spines['bottom'].set_visible(False)

    for tick in filter(lambda x: x >= 0, plt.get_yticks()):
        plt.axvline(x=tick, linestyle='dashed', alpha=0.4, color='#eeeeee',
                    zorder=1)

    return plt


def get_all_time_by_month_stat(data):
    df = prepare_dataframe(data)

    g = df.groupby(pd.Grouper(freq="M"))
    plt = g.sum().plot(figsize=(15, 8), color='#86bf91', fontsize=12,
                       title='for all time', grid=True, legend=False,
                       kind='line')

    plt.set_ylabel("Amount", labelpad=20, weight='bold', size=12)
    plt.set_xlabel("Month", labelpad=20, weight='bold', size=12)

    plt.tick_params(axis="both", which="both", bottom=False, top=False,
                    labelbottom=True, left=False, right=False, labelleft=True)

    return post_generate(plt)


def get_all_time_category_stat(data):
    df = prepare_dataframe(data)

    g = df.groupby(by='category')
    sm = g.sum().sort_values(by='amount', ascending=True)
    plt = sm.plot(figsize=(8, 10), color='#86bf91', fontsize=12, zorder=2,
                  width=0.85, title='categories for all time', kind='barh',
                  legend=False, )

    plt.set_ylabel("Category", labelpad=20, weight='bold', size=12)
    plt.set_xlabel("Amount", labelpad=20, weight='bold', size=12)

    for y, x in enumerate(sm['amount']):
        plt.annotate("%.2f" % x, xy=(x, y), va='center')

    return post_generate(plt)
